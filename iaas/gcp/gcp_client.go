package gcp

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	errwrap "github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"github.com/pivotal-cf/cliaas/iaas"
)

type GoogleComputeClient interface {
	List(project string, zone string) (*compute.InstanceList, error)
	DiskList(project string, zone string) (*compute.DiskList, error)
	Delete(project string, zone string, instanceName string) (*compute.Operation, error)
	Insert(project string, zone string, instance *compute.Instance) (*compute.Operation, error)
	ImageInsert(project string, image *compute.Image, timeout time.Duration) (*compute.Operation, error)
	Stop(project string, zone string, instanceName string) (*compute.Operation, error)
}

type ClientAPI interface {
	CreateVM(instance compute.Instance) error
	DeleteVM(instanceName string) error
	GetVMInfo(filter Filter) (*compute.Instance, error)
	Disk(filter Filter) (*compute.Disk, error)
	StopVM(instanceName string) error
	CreateImage(tarball string, diskSizeGB int64) (string, error)
	WaitForStatus(vmName string, desiredStatus string) error
}

type Client struct {
	projectName  string
	zoneName     string
	googleClient GoogleComputeClient
	timeout      time.Duration
}

//NewDefaultGoogleComputeClient -- builds a gcp client which connects to your gcp using `GOOGLE_APPLICATION_CREDENTIALS`
func NewDefaultGoogleComputeClient(credpath string) (GoogleComputeClient, error) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credpath)
	if err != nil {
		return nil, errwrap.Wrap(err, "couldnt set credentials ENV Var")
	}

	ctx := context.Background()
	hc, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		return nil, errwrap.Wrap(err, "we have a DefaultClient error")
	}

	c, err := compute.New(hc)
	if err != nil {
		return nil, errwrap.Wrap(err, "we have a compute.New error")
	}
	return &googleComputeClientWrapper{
		instanceService: c.Instances,
		disksService:    c.Disks,
		imageService:    c.Images,
		ctx:             ctx,
	}, nil
}

func NewClient(configs ...func(*Client) error) (*Client, error) {
	gcpClient := new(Client)
	gcpClient.timeout = 5 * time.Minute

	for _, cfg := range configs {
		err := cfg(gcpClient)
		if err != nil {
			return nil, errwrap.Wrap(err, "new GCP Client config loading error")
		}
	}

	if gcpClient.googleClient == nil {
		return nil, fmt.Errorf("You have an incomplete GCPClientAPI.googleClient")
	}

	if gcpClient.zoneName == "" {
		return nil, fmt.Errorf("You have an incomplete GCPClientAPI.zoneName")
	}

	if gcpClient.projectName == "" {
		return nil, fmt.Errorf("You have an incomplete GCPClientAPI.projectName")
	}
	return gcpClient, nil
}

/* Cliaas Client Interface */
func (c *Client) Delete(identifier string) error {
	return c.DeleteVM(identifier)
}

func (c *Client) Replace(identifier string, sourceImageTarballURL string, diskSizeGB int64) error {
	vmInstance, err := c.GetVMInfo(Filter{
		NameRegexString: identifier + "*",
	})
	if err != nil {
		return errwrap.Wrap(err, "getvminfo failed")
	}

	fmt.Printf("Stopping VM '%s'\n", vmInstance.Name)
	err = c.StopVM(vmInstance.Name)
	if err != nil {
		return errwrap.Wrap(err, "stopvm failed")
	}

	fmt.Printf("Waiting for  VM '%s' to terminate\n", vmInstance.Name)
	err = c.WaitForStatus(vmInstance.Name, InstanceTerminated)
	if err != nil {
		return errwrap.Wrap(err, "waitforstatus after stopvm failed")
	}

	fmt.Printf("Creating image with image '%s' and disk size '%d'GB\n", sourceImageTarballURL, diskSizeGB)
	sourceImage, err := c.CreateImage(sourceImageTarballURL, diskSizeGB)
	if err != nil {
		return errwrap.Wrap(err, "could not create new disk image")
	}

	fmt.Printf("Creating VM '%s' using image '%s', and disk size '%d'GB\n", fmt.Sprintf("%s-%s", identifier, time.Now().Format("2006-01-02-15-04-05")), sourceImage, diskSizeGB)
	newInstance := createGCPInstanceFromExisting(vmInstance, sourceImage, diskSizeGB, fmt.Sprintf("%s-%s", identifier, time.Now().Format("2006-01-02-15-04-05")))
	err = c.CreateVM(*newInstance)
	if err != nil {
		return errwrap.Wrap(err, "CreateVM call failed")
	}

	fmt.Printf("Waiting for  VM '%s' to start up\n", newInstance.Name)
	return c.WaitForStatus(newInstance.Name, InstanceRunning)
}

func ConfigTimeout(value time.Duration) func(*Client) error {
	return func(gcpClient *Client) error {
		gcpClient.timeout = value * time.Second
		return nil
	}
}

func (s *Client) GetDisk(identifier string) (iaas.Disk, error) {
	disk, err := s.Disk(Filter{
		NameRegexString: identifier + "*",
	})
	if err != nil {
		return iaas.Disk{}, err
	}
	return iaas.Disk{
		SizeGB: int64(disk.SizeGb),
	}, nil
}

/* End Cliaas Client Interface */

func (s *Client) Disk(filter Filter) (*compute.Disk, error) {
	list, err := s.googleClient.DiskList(s.projectName, s.zoneName)
	if err != nil {
		return nil, errwrap.Wrap(err, "call DiskList on google client failed")
	}

	for _, item := range list.Items {
		var validName = regexp.MustCompile(filter.NameRegexString)
		nameMatch := validName.MatchString(item.Name)
		if nameMatch {
			return item, nil
		}
	}
	return nil, fmt.Errorf("No disk matches found")
}

func ConfigGoogleClient(value GoogleComputeClient) func(*Client) error {
	return func(gcpClient *Client) error {
		gcpClient.googleClient = value
		return nil
	}
}

func ConfigZoneName(value string) func(*Client) error {
	return func(gcpClient *Client) error {
		gcpClient.zoneName = value
		return nil
	}
}

func ConfigProjectName(value string) func(*Client) error {
	return func(gcpClient *Client) error {
		gcpClient.projectName = value
		return nil
	}
}

func (s *Client) CreateImage(tarball string, diskSizeGB int64) (string, error) {
	imageName := fmt.Sprintf("opsman-disk-%v", time.Now().Format("2006-01-02-15-04-05"))
	fmt.Printf("Creating image '%s' with disk size '%d'GB \n", imageName, diskSizeGB)
	_, err := s.googleClient.ImageInsert(s.projectName, &compute.Image{
		Name:       imageName,
		DiskSizeGb: diskSizeGB,
		RawDisk: &compute.ImageRawDisk{
			Source: fmt.Sprintf("http://storage.googleapis.com/%v", tarball),
		},
	}, s.timeout)
	if err != nil {
		return "", err
	}

	sourceImage := fmt.Sprintf("projects/%s/global/images/%s", s.projectName, imageName)
	return sourceImage, nil
}

func (s *Client) CreateVM(instance compute.Instance) error {
	fmt.Printf("Creating VM '%s'\n", instance.Name)
	operation, err := s.googleClient.Insert(s.projectName, s.zoneName, &instance)
	if err != nil {
		return errwrap.Wrap(err, "call to googleclient.Insert yielded error")
	}

	if operation.Error != nil {
		return errors.New("unexpected errors from operation response from google client")
	}

	return nil
}

func (s *Client) DeleteVM(instanceName string) error {
	operation, err := s.googleClient.Delete(s.projectName, s.zoneName, instanceName)
	if err != nil {
		return errwrap.Wrap(err, "call to googleclient.Delete yielded error")
	}

	if operation.Error != nil {
		return errors.New("unexpected errors from operation response from google client")
	}

	return nil
}

//StopVM - will try to stop the VM with the given name
func (s *Client) StopVM(instanceName string) error {
	operation, err := s.googleClient.Stop(s.projectName, s.zoneName, instanceName)
	if err != nil {
		return errwrap.Wrap(err, "call to googleclient.Stop yielded error")
	}

	if operation.Error != nil {
		return errors.New("unexpected errors from operation response from google client")
	}

	return nil
}

//GetVMInfo - gets the information on the first VM to match the given filter argument
// currently filter will only do a regex on teh tag||name regex fields against
// the List's result set
func (s *Client) GetVMInfo(filter Filter) (*compute.Instance, error) {
	return s.getVMInfo(filter, InstanceRunning)
}

func (s *Client) getVMInfo(filter Filter, status string) (*compute.Instance, error) {
	list, err := s.googleClient.List(s.projectName, s.zoneName)
	if err != nil {
		return nil, errwrap.Wrap(err, "call List on google client failed")
	}

	for _, item := range list.Items {
		var validID = regexp.MustCompile(filter.TagRegexString)
		var validName = regexp.MustCompile(filter.NameRegexString)
		taglist := strings.Join(item.Tags.Items, " ")
		tagMatch := validID.MatchString(taglist)
		nameMatch := validName.MatchString(item.Name)

		if tagMatch &&
			nameMatch &&
			(status == InstanceAll || item.Status == InstanceRunning) {
			return item, nil
		}
	}
	return nil, fmt.Errorf("No instance matches found")
}

func (s *Client) WaitForStatus(vmName string, desiredStatus string) error {
	errChannel := make(chan error)
	go func() {
		for {
			vmInfo, err := s.getVMInfo(Filter{NameRegexString: vmName}, InstanceAll)
			if err != nil {
				errChannel <- errwrap.Wrap(err, "GetVMInfo call failed")
				return
			}

			if vmInfo.Status == desiredStatus {
				errChannel <- nil
				return
			}
		}
	}()
	select {
	case res := <-errChannel:
		return res
	case <-time.After(s.timeout):
		return errors.New("polling for status timed out")
	}
}

type googleComputeClientWrapper struct {
	imageService    *compute.ImagesService
	instanceService *compute.InstancesService
	disksService    *compute.DisksService
	ctx             context.Context
}

func (s *googleComputeClientWrapper) List(project string, zone string) (*compute.InstanceList, error) {
	return s.instanceService.List(project, zone).Context(s.ctx).Do()
}

func (s *googleComputeClientWrapper) Delete(project string, zone string, instance string) (*compute.Operation, error) {
	return s.instanceService.Delete(project, zone, instance).Context(s.ctx).Do()
}

func (s *googleComputeClientWrapper) Stop(project string, zone string, instance string) (*compute.Operation, error) {
	vmInstance, err := s.instanceService.Get(project, zone, instance).Context(s.ctx).Do()
	if err != nil {
		return nil, errwrap.Wrap(err, "failed getting vm instance")
	}

	if len(vmInstance.NetworkInterfaces) > 0 && len(vmInstance.NetworkInterfaces[0].AccessConfigs) > 0 {
		accessConfigName := vmInstance.NetworkInterfaces[0].AccessConfigs[0].Name
		nicName := vmInstance.NetworkInterfaces[0].Name
		operation, err := s.instanceService.DeleteAccessConfig(project, zone, instance, accessConfigName, nicName).Context(s.ctx).Do()
		if err != nil {
			return operation, errwrap.Wrap(err, "could not delete access config")
		}

	}

	return s.instanceService.Stop(project, zone, instance).Context(s.ctx).Do()
}

func (s *googleComputeClientWrapper) Insert(project string, zone string, instance *compute.Instance) (*compute.Operation, error) {
	return s.instanceService.Insert(project, zone, instance).Context(s.ctx).Do()
}

func (s *googleComputeClientWrapper) DiskList(project string, zone string) (*compute.DiskList, error) {
	return s.disksService.List(project, zone).Context(s.ctx).Do()
}

func (s *googleComputeClientWrapper) ImageInsert(project string, image *compute.Image, timeout time.Duration) (*compute.Operation, error) {
	operation, err := s.imageService.Insert(project, image).Context(s.ctx).Do()
	if err != nil {
		return operation, errwrap.Wrap(err, "disk image insert failed")
	}

	errChannel := make(chan error)
	go func() {
		for {
			image, err := s.imageService.Get(project, image.Name).Context(s.ctx).Do()
			if err != nil {
				errChannel <- errwrap.Wrap(err, "image get failed")
			}

			if image != nil && image.Status == ImageReady {
				errChannel <- nil
			}

			if image != nil && image.Status == ImageFailed {
				errChannel <- errors.New("image creation failed")
			}
		}
	}()
	select {
	case res := <-errChannel:
		return operation, res
	case <-time.After(timeout):
		return nil, errors.New("polling for status timed out")
	}
	return nil, nil
}

func createGCPInstanceFromExisting(vmInstance *compute.Instance, sourceImage string, diskSizeGB int64, name string) *compute.Instance {
	newInstance := &compute.Instance{
		NetworkInterfaces: vmInstance.NetworkInterfaces,
		MachineType:       vmInstance.MachineType,
		Name:              name,
		Tags: &compute.Tags{
			Items: vmInstance.Tags.Items,
		},
		Disks: []*compute.AttachedDisk{
			&compute.AttachedDisk{
				Boot: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: sourceImage,
					DiskSizeGb:  diskSizeGB,
				},
			},
		},
	}
	newInstance.NetworkInterfaces[0].NetworkIP = ""
	return newInstance
}

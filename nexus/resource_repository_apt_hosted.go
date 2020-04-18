package nexus

import (
	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceRepositoryAptHosted() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryCreate,
		Read:   resourceRepositoryRead,
		Update: resourceRepositoryUpdate,
		Delete: resourceRepositoryDelete,
		Exists: resourceRepositoryExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			// "format": {
			// 	Description:  "Repository format",
			// 	ForceNew:     true,
			// 	Required:     true,
			// 	Type:         schema.TypeString,
			// 	ValidateFunc: validation.StringInSlice([]string{"apt", "bower", "docker", "maven2", "nuget", "pypi"}, false),
			// },
			"name": {
				Description: "A unique identifier for this repository",
				Required:    true,
				Type:        schema.TypeString,
			},
			"online": {
				Default:     true,
				Description: "Whether this repository accepts incoming requests",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			// "type": {
			// 	Description:  "Repository type",
			// 	ForceNew:     true,
			// 	Type:         schema.TypeString,
			// 	Required:     true,
			// 	ValidateFunc: validation.StringInSlice([]string{"group", "hosted", "proxy"}, false),
			// },
			"apt": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"distribution": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"apt_signing": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keypair": {
							Type:     schema.TypeString,
							Required: true,
						},
						"passphrase": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"cleanup": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_names": {
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
							Type:     schema.TypeSet,
						},
					},
				},
			},
			"storage": {
				DefaultFunc: repositoryStorageDefault,
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blob_store_name": {
							Default:     "default",
							Description: "Blob store used to store repository contents",
							Optional:    true,
							Type:        schema.TypeString,
						},
						"strict_content_type_validation": {
							Default:     true,
							Description: "Whether to validate uploaded content's MIME type appropriate for the repository format",
							Optional:    true,
							Type:        schema.TypeBool,
						},
						"write_policy": {
							Description: "Controls if deployments of and updates to assets are allowed",
							Optional:    true,
							Type:        schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"ALLOW",
								"ALLOW_ONCE",
								"DENY",
							}, false),
						},
					},
				},
			},
		},
	}
}
func getRepositoryFromAptHostedResourceData(d *schema.ResourceData) nexus.Repository {
	repo := nexus.Repository{
		Name:   d.Get("name").(string),
		Online: d.Get("online").(bool),
		Type:   d.Get("type").(string),
		Format: d.Get("format").(string),
	}
	if _, ok := d.GetOk("apt"); ok {
		aptList := d.Get("apt").([]interface{})
		aptConfig := aptList[0].(map[string]interface{})
		repo.RepositoryApt = &nexus.RepositoryApt{
			Distribution: aptConfig["distribution"].(string),
		}
	}
	if _, ok := d.GetOk("apt_signing"); ok {
		aptSigningList := d.Get("apt_signing").([]interface{})
		aptSigningConfig := aptSigningList[0].(map[string]interface{})
		repo.RepositoryAptSigning = &nexus.RepositoryAptSigning{
			Keypair:    aptSigningConfig["keypair"].(string),
			Passphrase: aptSigningConfig["passphrase"].(string),
		}
	}
	if _, ok := d.GetOk("cleanup"); ok {
		cleanupList := d.Get("cleanup").([]interface{})
		cleanupConfig := cleanupList[0].(map[string]interface{})
		repoCleanupPolicyNames := cleanupConfig["policy_names"].(*schema.Set)
		cleanupPolicyNames := make([]string, repoCleanupPolicyNames.Len())
		for _, v := range repoCleanupPolicyNames.List() {
			cleanupPolicyNames = append(cleanupPolicyNames, v.(string))
		}
		repo.RepositoryCleanup = &nexus.RepositoryCleanup{
			PolicyNames: cleanupPolicyNames,
		}
	}
	if _, ok := d.GetOk("storage"); ok {
		storageList := d.Get("storage").([]interface{})
		storageConfig := storageList[0].(map[string]interface{})
		repo.RepositoryStorage = &nexus.RepositoryStorage{
			BlobStoreName:               storageConfig["blob_store_name"].(string),
			StrictContentTypeValidation: storageConfig["strict_content_type_validation"].(bool),
			WritePolicy:                 storageConfig["write_policy"].(string),
		}
	}
	return repo
}
func setRepositoryAptHostedToResourceData(repo *nexus.Repository, d *schema.ResourceData) error {
	d.SetId(repo.Name)
	d.Set("name", repo.Name)
	d.Set("online", repo.Online)
	d.Set("format", "apt")
	d.Set("type", "hosted")
	if repo.RepositoryApt != nil {
		if err := d.Set("apt", flattenRepositoryApt(repo.RepositoryApt)); err != nil {
			return err
		}
	}
	if repo.RepositoryAptSigning != nil {
		if err := d.Set("apt_signing", flattenRepositoryAptSigning(repo.RepositoryAptSigning)); err != nil {
			return err
		}
	}

	if err := d.Set("storage", flattenRepositoryStorage(repo.RepositoryStorage)); err != nil {
		return err
	}
	return nil
}

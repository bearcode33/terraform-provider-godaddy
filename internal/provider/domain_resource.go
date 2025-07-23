package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bearcode33/terraform-provider-godaddy/internal/godaddy"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &DomainResource{}
var _ resource.ResourceWithImportState = &DomainResource{}

func NewDomainResource() resource.Resource {
	return &DomainResource{}
}

type DomainResource struct {
	client *godaddy.Client
}

type DomainResourceModel struct {
	Domain              types.String `tfsdk:"domain"`
	Status              types.String `tfsdk:"status"`
	Expires             types.String `tfsdk:"expires"`
	ExpirationProtected types.Bool   `tfsdk:"expiration_protected"`
	HoldRegistrar       types.Bool   `tfsdk:"hold_registrar"`
	Locked              types.Bool   `tfsdk:"locked"`
	Privacy             types.Bool   `tfsdk:"privacy"`
	RenewAuto           types.Bool   `tfsdk:"renew_auto"`
	TransferProtected   types.Bool   `tfsdk:"transfer_protected"`
	Nameservers         types.List   `tfsdk:"nameservers"`
	ContactAdmin        types.Object `tfsdk:"contact_admin"`
	ContactBilling      types.Object `tfsdk:"contact_billing"`
	ContactRegistrant   types.Object `tfsdk:"contact_registrant"`
	ContactTech         types.Object `tfsdk:"contact_tech"`
}

func (r *DomainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (r *DomainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a GoDaddy domain and its settings.",
		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				MarkdownDescription: "The domain name.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The current status of the domain.",
				Computed:            true,
			},
			"expires": schema.StringAttribute{
				MarkdownDescription: "The expiration date of the domain.",
				Computed:            true,
			},
			"expiration_protected": schema.BoolAttribute{
				MarkdownDescription: "Whether the domain is protected from expiration.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"hold_registrar": schema.BoolAttribute{
				MarkdownDescription: "Whether the domain has a registrar hold.",
				Computed:            true,
			},
			"locked": schema.BoolAttribute{
				MarkdownDescription: "Whether the domain is locked to prevent transfers.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"privacy": schema.BoolAttribute{
				MarkdownDescription: "Whether WHOIS privacy is enabled.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"renew_auto": schema.BoolAttribute{
				MarkdownDescription: "Whether the domain is set to auto-renew.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"transfer_protected": schema.BoolAttribute{
				MarkdownDescription: "Whether the domain is protected from transfers.",
				Computed:            true,
			},
			"nameservers": schema.ListAttribute{
				MarkdownDescription: "The nameservers for the domain.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"contact_admin": schema.SingleNestedAttribute{
				MarkdownDescription: "Administrative contact information.",
				Optional:            true,
				Computed:            true,
				Attributes:          contactAttributes(),
			},
			"contact_billing": schema.SingleNestedAttribute{
				MarkdownDescription: "Billing contact information.",
				Optional:            true,
				Computed:            true,
				Attributes:          contactAttributes(),
			},
			"contact_registrant": schema.SingleNestedAttribute{
				MarkdownDescription: "Registrant contact information.",
				Optional:            true,
				Computed:            true,
				Attributes:          contactAttributes(),
			},
			"contact_tech": schema.SingleNestedAttribute{
				MarkdownDescription: "Technical contact information.",
				Optional:            true,
				Computed:            true,
				Attributes:          contactAttributes(),
			},
		},
	}
}

func contactAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name_first": schema.StringAttribute{
			MarkdownDescription: "First name.",
			Required:            true,
		},
		"name_middle": schema.StringAttribute{
			MarkdownDescription: "Middle name.",
			Optional:            true,
		},
		"name_last": schema.StringAttribute{
			MarkdownDescription: "Last name.",
			Required:            true,
		},
		"organization": schema.StringAttribute{
			MarkdownDescription: "Organization name.",
			Optional:            true,
		},
		"job_title": schema.StringAttribute{
			MarkdownDescription: "Job title.",
			Optional:            true,
		},
		"email": schema.StringAttribute{
			MarkdownDescription: "Email address.",
			Required:            true,
		},
		"phone": schema.StringAttribute{
			MarkdownDescription: "Phone number.",
			Required:            true,
		},
		"fax": schema.StringAttribute{
			MarkdownDescription: "Fax number.",
			Optional:            true,
		},
		"address1": schema.StringAttribute{
			MarkdownDescription: "Address line 1.",
			Required:            true,
		},
		"address2": schema.StringAttribute{
			MarkdownDescription: "Address line 2.",
			Optional:            true,
		},
		"city": schema.StringAttribute{
			MarkdownDescription: "City.",
			Required:            true,
		},
		"state": schema.StringAttribute{
			MarkdownDescription: "State or province.",
			Required:            true,
		},
		"postal_code": schema.StringAttribute{
			MarkdownDescription: "Postal code.",
			Required:            true,
		},
		"country": schema.StringAttribute{
			MarkdownDescription: "Country code (2-letter ISO).",
			Required:            true,
		},
	}
}

func (r *DomainResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*godaddy.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *godaddy.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *DomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DomainResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For existing domains, we'll just read the current state
	domain, err := r.client.GetDomain(ctx, data.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Domain",
			fmt.Sprintf("Could not read domain %s: %s", data.Domain.ValueString(), err),
		)
		return
	}

	// Update the model with the domain data
	r.updateModelFromDomain(&data, domain)

	// Apply any configuration changes
	if err := r.applyDomainSettings(ctx, &data, domain); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Domain Settings",
			fmt.Sprintf("Could not update domain settings for %s: %s", data.Domain.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "created domain resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DomainResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain, err := r.client.GetDomain(ctx, data.Domain.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Domain",
			fmt.Sprintf("Could not read domain %s: %s", data.Domain.ValueString(), err),
		)
		return
	}

	r.updateModelFromDomain(&data, domain)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DomainResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain, err := r.client.GetDomain(ctx, data.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Domain",
			fmt.Sprintf("Could not read domain %s: %s", data.Domain.ValueString(), err),
		)
		return
	}

	if err := r.applyDomainSettings(ctx, &data, domain); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Domain Settings",
			fmt.Sprintf("Could not update domain settings for %s: %s", data.Domain.ValueString(), err),
		)
		return
	}

	// Read back the updated domain
	domain, err = r.client.GetDomain(ctx, data.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Domain",
			fmt.Sprintf("Could not read updated domain %s: %s", data.Domain.ValueString(), err),
		)
		return
	}

	r.updateModelFromDomain(&data, domain)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DomainResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Note: GoDaddy API doesn't support deleting domains directly
	// The domain will remain in the account but Terraform will no longer manage it
	tflog.Warn(ctx, "Domain resource removed from Terraform state but remains in GoDaddy account",
		map[string]interface{}{"domain": data.Domain.ValueString()})
}

func (r *DomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("domain"), req, resp)
}

func (r *DomainResource) updateModelFromDomain(model *DomainResourceModel, domain *godaddy.DomainDetail) {
	model.Status = types.StringValue(domain.Status)

	if domain.Expires != nil {
		model.Expires = types.StringValue(domain.Expires.Format(time.RFC3339))
	}

	model.ExpirationProtected = types.BoolValue(domain.ExpirationProtected)
	model.HoldRegistrar = types.BoolValue(domain.HoldRegistrar)
	model.Locked = types.BoolValue(domain.Locked)
	model.Privacy = types.BoolValue(domain.Privacy)
	model.RenewAuto = types.BoolValue(domain.RenewAuto)
	model.TransferProtected = types.BoolValue(domain.TransferProtected)

	// Convert nameservers
	if len(domain.Nameservers) > 0 {
		nameservers := make([]types.String, len(domain.Nameservers))
		for i, ns := range domain.Nameservers {
			nameservers[i] = types.StringValue(ns)
		}
		model.Nameservers = types.ListValueMust(types.StringType, convertStringsToValues(nameservers))
	}

	// Convert contacts
	model.ContactAdmin = contactToObject(domain.ContactAdmin)
	model.ContactBilling = contactToObject(domain.ContactBilling)
	model.ContactRegistrant = contactToObject(domain.ContactRegistrant)
	model.ContactTech = contactToObject(domain.ContactTech)
}

func (r *DomainResource) applyDomainSettings(ctx context.Context, model *DomainResourceModel, currentDomain *godaddy.DomainDetail) error {
	update := godaddy.DomainUpdate{}
	needsUpdate := false

	// Check if locked state needs update
	if !model.Locked.IsNull() && model.Locked.ValueBool() != currentDomain.Locked {
		locked := model.Locked.ValueBool()
		update.Locked = &locked
		needsUpdate = true
	}

	// Check if renew_auto needs update
	if !model.RenewAuto.IsNull() && model.RenewAuto.ValueBool() != currentDomain.RenewAuto {
		renewAuto := model.RenewAuto.ValueBool()
		update.RenewAuto = &renewAuto
		needsUpdate = true
	}

	// Check if nameservers need update
	if !model.Nameservers.IsNull() {
		var nameservers []string
		model.Nameservers.ElementsAs(ctx, &nameservers, false)

		if !stringSlicesEqual(nameservers, currentDomain.Nameservers) {
			update.NameServers = nameservers
			needsUpdate = true
		}
	}

	if needsUpdate {
		if err := r.client.UpdateDomain(ctx, model.Domain.ValueString(), update); err != nil {
			return err
		}
	}

	// Update contacts if provided
	if err := r.updateContacts(ctx, model, currentDomain); err != nil {
		return err
	}

	return nil
}

func (r *DomainResource) updateContacts(ctx context.Context, model *DomainResourceModel, currentDomain *godaddy.DomainDetail) error {
	contacts := godaddy.DomainContacts{}
	needsUpdate := false

	if !model.ContactAdmin.IsNull() {
		contact := objectToContact(model.ContactAdmin)
		if contact != nil && !contactsEqual(*contact, currentDomain.ContactAdmin) {
			contacts.ContactAdmin = contact
			needsUpdate = true
		}
	}

	if !model.ContactBilling.IsNull() {
		contact := objectToContact(model.ContactBilling)
		if contact != nil && !contactsEqual(*contact, currentDomain.ContactBilling) {
			contacts.ContactBilling = contact
			needsUpdate = true
		}
	}

	if !model.ContactRegistrant.IsNull() {
		contact := objectToContact(model.ContactRegistrant)
		if contact != nil && !contactsEqual(*contact, currentDomain.ContactRegistrant) {
			contacts.ContactRegistrant = contact
			needsUpdate = true
		}
	}

	if !model.ContactTech.IsNull() {
		contact := objectToContact(model.ContactTech)
		if contact != nil && !contactsEqual(*contact, currentDomain.ContactTech) {
			contacts.ContactTech = contact
			needsUpdate = true
		}
	}

	if needsUpdate {
		return r.client.UpdateDomainContacts(ctx, model.Domain.ValueString(), contacts)
	}

	return nil
}

func convertStringsToValues(strings []types.String) []attr.Value {
	values := make([]attr.Value, len(strings))
	for i, s := range strings {
		values[i] = s
	}
	return values
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func contactToObject(contact godaddy.DomainContact) types.Object {
	attributes := map[string]attr.Value{
		"name_first":   types.StringValue(contact.NameFirst),
		"name_middle":  types.StringValue(contact.NameMiddle),
		"name_last":    types.StringValue(contact.NameLast),
		"organization": types.StringValue(contact.Organization),
		"job_title":    types.StringValue(contact.JobTitle),
		"email":        types.StringValue(contact.Email),
		"phone":        types.StringValue(contact.Phone),
		"fax":          types.StringValue(contact.Fax),
		"address1":     types.StringValue(contact.AddressMailing.Address1),
		"address2":     types.StringValue(contact.AddressMailing.Address2),
		"city":         types.StringValue(contact.AddressMailing.City),
		"state":        types.StringValue(contact.AddressMailing.State),
		"postal_code":  types.StringValue(contact.AddressMailing.PostalCode),
		"country":      types.StringValue(contact.AddressMailing.Country),
	}

	return types.ObjectValueMust(contactAttributeTypes(), attributes)
}

func objectToContact(obj types.Object) *godaddy.DomainContact {
	if obj.IsNull() {
		return nil
	}

	attrs := obj.Attributes()
	return &godaddy.DomainContact{
		NameFirst:    attrs["name_first"].(types.String).ValueString(),
		NameMiddle:   attrs["name_middle"].(types.String).ValueString(),
		NameLast:     attrs["name_last"].(types.String).ValueString(),
		Organization: attrs["organization"].(types.String).ValueString(),
		JobTitle:     attrs["job_title"].(types.String).ValueString(),
		Email:        attrs["email"].(types.String).ValueString(),
		Phone:        attrs["phone"].(types.String).ValueString(),
		Fax:          attrs["fax"].(types.String).ValueString(),
		AddressMailing: godaddy.DomainAddress{
			Address1:   attrs["address1"].(types.String).ValueString(),
			Address2:   attrs["address2"].(types.String).ValueString(),
			City:       attrs["city"].(types.String).ValueString(),
			State:      attrs["state"].(types.String).ValueString(),
			PostalCode: attrs["postal_code"].(types.String).ValueString(),
			Country:    attrs["country"].(types.String).ValueString(),
		},
	}
}

func contactAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name_first":   types.StringType,
		"name_middle":  types.StringType,
		"name_last":    types.StringType,
		"organization": types.StringType,
		"job_title":    types.StringType,
		"email":        types.StringType,
		"phone":        types.StringType,
		"fax":          types.StringType,
		"address1":     types.StringType,
		"address2":     types.StringType,
		"city":         types.StringType,
		"state":        types.StringType,
		"postal_code":  types.StringType,
		"country":      types.StringType,
	}
}

func contactsEqual(a, b godaddy.DomainContact) bool {
	return a.NameFirst == b.NameFirst &&
		a.NameMiddle == b.NameMiddle &&
		a.NameLast == b.NameLast &&
		a.Organization == b.Organization &&
		a.JobTitle == b.JobTitle &&
		a.Email == b.Email &&
		a.Phone == b.Phone &&
		a.Fax == b.Fax &&
		a.AddressMailing.Address1 == b.AddressMailing.Address1 &&
		a.AddressMailing.Address2 == b.AddressMailing.Address2 &&
		a.AddressMailing.City == b.AddressMailing.City &&
		a.AddressMailing.State == b.AddressMailing.State &&
		a.AddressMailing.PostalCode == b.AddressMailing.PostalCode &&
		a.AddressMailing.Country == b.AddressMailing.Country
}

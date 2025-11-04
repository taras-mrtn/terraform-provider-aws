// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pinpointsmsvoicev2

import (
	"context"
	"fmt"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"
	awstypes "github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkid "github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/fwdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	fwflex "github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	fwtypes "github.com/hashicorp/terraform-provider-aws/internal/framework/types"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @FrameworkResource("aws_pinpointsmsvoicev2_sender_id", name="Sender ID")
// @Tags(identifierAttribute="arn")
func newSenderIdResource(context.Context) (resource.ResourceWithConfigure, error) {
	r := &senderIdResource{}

	return r, nil
}

type senderIdResource struct {
	framework.ResourceWithModel[senderIdResourceModel]
	framework.WithImportByID
}

func (r *senderIdResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			names.AttrARN: framework.ARNAttributeComputedOnly(),
			"deletion_protection_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			names.AttrID: framework.IDAttribute(),
			"iso_country_code": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexache.MustCompile(`^[A-Z]{2}$`), "must be in ISO 3166-1 alpha-2 format"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"message_types": schema.SetAttribute{
				CustomType:  fwtypes.NewSetTypeOf[fwtypes.StringEnum[awstypes.MessageType]](ctx),
				Optional:    true,
				ElementType: fwtypes.StringEnumType[awstypes.MessageType](),
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.SizeAtMost(2),
				},
				PlanModifiers: []planmodifier.Set{},
			},
			"monthly_leasing_price": schema.StringAttribute{
				Computed: true,
			},
			"registered": schema.BoolAttribute{
				Computed: true,
			},
			"sender_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexache.MustCompile(`^[A-Za-z0-9_-]{1,11}$`), "must be between 1 and 11 characters long and contain only letters, numbers, underscores, and dashes"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			names.AttrTags:    tftags.TagsAttribute(),
			names.AttrTagsAll: tftags.TagsAttributeComputedOnly(),
		},
	}
}

func (r *senderIdResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data senderIdResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	conn := r.Meta().PinpointSMSVoiceV2Client(ctx)

	input := &pinpointsmsvoicev2.RequestSenderIdInput{}
	response.Diagnostics.Append(fwflex.Expand(ctx, data, input)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Additional fields.
	input.ClientToken = aws.String(sdkid.UniqueId())
	input.Tags = getTagsIn(ctx)

	output, err := conn.RequestSenderId(ctx, input)

	if err != nil {
		response.Diagnostics.AddError("requesting End User Messaging SMS Sender ID", err.Error())

		return
	}

	// Set values for unknowns.
	data.SenderIdARN = fwflex.StringToFramework(ctx, output.SenderIdArn)
	data.setID()

	// Read back to get all computed attributes
	out, err := FindSenderIdByARN(ctx, conn, data.ID.ValueString())

	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("reading End User Messaging SMS Sender ID (%s)", data.ID.ValueString()), err.Error())

		return
	}

	response.Diagnostics.Append(fwflex.Flatten(ctx, out, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, data)...)
}

func (r *senderIdResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data senderIdResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	if err := data.InitFromID(); err != nil {
		response.Diagnostics.AddError("parsing resource ID", err.Error())

		return
	}

	conn := r.Meta().PinpointSMSVoiceV2Client(ctx)

	out, err := FindSenderIdByARN(ctx, conn, data.ID.ValueString())

	if tfresource.NotFound(err) {
		response.Diagnostics.Append(fwdiag.NewResourceNotFoundWarningDiagnostic(err))
		response.State.RemoveResource(ctx)

		return
	}

	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("reading End User Messaging SMS Sender ID (%s)", data.ID.ValueString()), err.Error())

		return
	}

	// Set attributes for import.
	response.Diagnostics.Append(fwflex.Flatten(ctx, out, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *senderIdResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var old, new senderIdResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &new)...)
	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(request.State.Get(ctx, &old)...)
	if response.Diagnostics.HasError() {
		return
	}

	conn := r.Meta().PinpointSMSVoiceV2Client(ctx)

	if !new.DeletionProtectionEnabled.Equal(old.DeletionProtectionEnabled) {
		input := &pinpointsmsvoicev2.UpdateSenderIdInput{
			SenderId:                  fwflex.StringFromFramework(ctx, new.SenderID),
			IsoCountryCode:            fwflex.StringFromFramework(ctx, new.ISOCountryCode),
			DeletionProtectionEnabled: fwflex.BoolFromFramework(ctx, new.DeletionProtectionEnabled),
		}

		_, err := conn.UpdateSenderId(ctx, input)

		if err != nil {
			response.Diagnostics.AddError(fmt.Sprintf("updating End User Messaging SMS Sender ID (%s)", new.ID.ValueString()), err.Error())

			return
		}
	}

	response.Diagnostics.Append(response.State.Set(ctx, &new)...)
}

func (r *senderIdResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data senderIdResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	conn := r.Meta().PinpointSMSVoiceV2Client(ctx)

	_, err := conn.ReleaseSenderId(ctx, &pinpointsmsvoicev2.ReleaseSenderIdInput{
		SenderId:       data.SenderID.ValueStringPointer(),
		IsoCountryCode: data.ISOCountryCode.ValueStringPointer(),
	})

	if errs.IsA[*awstypes.ResourceNotFoundException](err) {
		return
	}

	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("releasing End User Messaging SMS Sender ID (%s)", data.ID.ValueString()), err.Error())

		return
	}
}

type senderIdResourceModel struct {
	framework.WithRegionModel
	ID                        types.String                                       `tfsdk:"id"`
	SenderIdARN               types.String                                       `tfsdk:"arn"`
	SenderID                  types.String                                       `tfsdk:"sender_id"`
	ISOCountryCode            types.String                                       `tfsdk:"iso_country_code"`
	MessageTypes              fwtypes.SetOfStringEnum[awstypes.MessageType]      `tfsdk:"message_types"`
	MonthlyLeasingPrice       types.String                                       `tfsdk:"monthly_leasing_price"`
	Registered                types.Bool                                         `tfsdk:"registered"`
	DeletionProtectionEnabled types.Bool                                         `tfsdk:"deletion_protection_enabled"`
	Tags                      tftags.Map                                         `tfsdk:"tags"`
	TagsAll                   tftags.Map                                         `tfsdk:"tags_all"`
}

func (model *senderIdResourceModel) InitFromID() error {
	model.SenderIdARN = model.ID

	return nil
}

func (model *senderIdResourceModel) setID() {
	model.ID = model.SenderIdARN
}

func FindSenderIdByARN(ctx context.Context, conn *pinpointsmsvoicev2.Client, arn string) (*awstypes.SenderIdInformation, error) {
	input := &pinpointsmsvoicev2.DescribeSenderIdsInput{
		SenderIds: []awstypes.SenderIdAndCountry{
			{
				SenderIdArn: aws.String(arn),
			},
		},
	}

	return findSenderId(ctx, conn, input)
}

func findSenderId(ctx context.Context, conn *pinpointsmsvoicev2.Client, input *pinpointsmsvoicev2.DescribeSenderIdsInput) (*awstypes.SenderIdInformation, error) {
	output, err := findSenderIds(ctx, conn, input)

	if err != nil {
		return nil, err
	}

	return tfresource.AssertSingleValueResult(output)
}

func findSenderIds(ctx context.Context, conn *pinpointsmsvoicev2.Client, input *pinpointsmsvoicev2.DescribeSenderIdsInput) ([]awstypes.SenderIdInformation, error) {
	var output []awstypes.SenderIdInformation

	pages := pinpointsmsvoicev2.NewDescribeSenderIdsPaginator(conn, input)
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)

		if errs.IsA[*awstypes.ResourceNotFoundException](err) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: input,
			}
		}

		if err != nil {
			return nil, err
		}

		output = append(output, page.SenderIds...)
	}

	return output, nil
}

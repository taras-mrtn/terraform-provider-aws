---
subcategory: "End User Messaging SMS"
layout: "aws"
page_title: "AWS: aws_pinpointsmsvoicev2_sender_id"
description: |-
  Manages an AWS End User Messaging SMS Sender ID.
---

# Resource: aws_pinpointsmsvoicev2_sender_id

Manages an AWS End User Messaging SMS Sender ID.

Sender IDs are alphanumeric names that identify the sender of an SMS message. When you send an SMS message using a sender ID, the sender ID appears on recipients' devices instead of a phone number. Sender IDs provide better brand recognition compared to sending messages from a random phone number.

~> **NOTE:** Sender IDs are only supported in certain countries. Some countries require pre-registration of sender IDs with local telecommunications authorities before they can be used. Review the [AWS End User Messaging SMS documentation](https://docs.aws.amazon.com/sms-voice/latest/userguide/sender-id.html) for country-specific requirements.

~> **NOTE:** This resource may incur monthly charges. Review the `monthly_leasing_price` attribute to understand the cost for your selected country.

## Example Usage

### Basic Usage

```terraform
resource "aws_pinpointsmsvoicev2_sender_id" "example" {
  sender_id        = "MyCompany"
  iso_country_code = "US"
}
```

### With Message Types

```terraform
resource "aws_pinpointsmsvoicev2_sender_id" "example" {
  sender_id        = "MyCompany"
  iso_country_code = "GB"

  message_types = [
    "TRANSACTIONAL",
    "PROMOTIONAL",
  ]
}
```

### With Deletion Protection and Tags

```terraform
resource "aws_pinpointsmsvoicev2_sender_id" "example" {
  sender_id                   = "MyCompany"
  iso_country_code            = "IN"
  deletion_protection_enabled = true

  tags = {
    Environment = "production"
    Application = "notifications"
  }
}
```

## Argument Reference

This resource supports the following arguments:

* `sender_id` - (Required) The sender ID string to request. Must be between 1 and 11 characters long and contain only letters, numbers, underscores, and dashes. Changing this value will create a new resource.
* `iso_country_code` - (Required) The two-character code, in ISO 3166-1 alpha-2 format, for the country or region where you want to request the sender ID. Changing this value will create a new resource.
* `message_types` - (Optional) The type of messages that can be sent using the sender ID. Valid values are `TRANSACTIONAL` (for messages that are critical or time-sensitive) and `PROMOTIONAL` (for messages that aren't critical or time-sensitive). If not specified, both message types are allowed.
* `deletion_protection_enabled` - (Optional) By default this is set to `false`. When set to `true`, the sender ID can't be deleted.
* `tags` - (Optional) Key-value map of resource tags. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `id` - The ARN of the sender ID.
* `arn` - The ARN of the sender ID.
* `monthly_leasing_price` - The monthly price, in US dollars, to lease the sender ID.
* `registered` - Whether the sender ID is registered with the wireless carrier.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block).

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import sender IDs using the `arn`. For example:

```terraform
import {
  to = aws_pinpointsmsvoicev2_sender_id.example
  id = "arn:aws:sms-voice:us-east-1:123456789012:sender-id/12345678-1234-1234-1234-123456789012/MyCompany"
}
```

Using `terraform import`, import sender IDs using the `arn`. For example:

```console
% terraform import aws_pinpointsmsvoicev2_sender_id.example arn:aws:sms-voice:us-east-1:123456789012:sender-id/12345678-1234-1234-1234-123456789012/MyCompany
```

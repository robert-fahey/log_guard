variable "project_id" {
  description = "The ID of the project in which the resource belongs"
}

locals {
  base_filter = ["logName:\"/logs/error\""]
}

resource "google_logging_metric" "error_null_pointer" {
  project = var.project_id
  name = "error_null_pointer"
  filter = join(" AND ", concat(local.base_filter, tolist([
    "jsonPayload.code=\"ERROR_NULL_POINTER\""
  ])))
  metric_descriptor {
    value_type  = "INT64"
    unit        = "1"
    metric_kind = "DELTA"
  }
}

resource "google_logging_metric" "warn_data_format_mismatch" {
  project = var.project_id
  name = "warn_data_format_mismatch"
  filter = join(" AND ", concat(local.base_filter, tolist([
    "jsonPayload.code=\"WARN_DATA_FORMAT_MISMATCH\""
  ])))
  metric_descriptor {
    value_type  = "INT64"
    unit        = "1"
    metric_kind = "DELTA"
  }
}
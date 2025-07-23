variable "godaddy_api_key" {
  description = "GoDaddy API Key"
  type        = string
  sensitive   = true
}

variable "godaddy_api_secret" {
  description = "GoDaddy API Secret"
  type        = string
  sensitive   = true
}

variable "domain_name" {
  description = "The domain name to manage"
  type        = string
  default     = "example.com"
}

variable "environment" {
  description = "GoDaddy API environment (production or test)"
  type        = string
  default     = "production"

  validation {
    condition     = contains(["production", "test"], var.environment)
    error_message = "Environment must be either 'production' or 'test'."
  }
}
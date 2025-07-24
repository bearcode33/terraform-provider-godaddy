# Optional: Define variables for API credentials
# It's recommended to use environment variables instead

variable "godaddy_api_key" {
  description = "GoDaddy API Key"
  type        = string
  sensitive   = true
  default     = null
}

variable "godaddy_api_secret" {
  description = "GoDaddy API Secret"
  type        = string
  sensitive   = true
  default     = null
}

variable "godaddy_environment" {
  description = "GoDaddy API Environment (production or test)"
  type        = string
  default     = "production"
  
  validation {
    condition     = contains(["production", "test"], var.godaddy_environment)
    error_message = "Environment must be 'production' or 'test'."
  }
}
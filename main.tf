# adapted from https://www.sethvargo.com/configuring-cloud-run-with-terraform/
terraform {
  required_version = ">= 0.12"
  required_providers {
    # Cloud Run resources were not added until 3.3.0
    google = ">= 3.3"
  }
}

variable "db_password" {
  type = string
}
variable "google_client_id" {
  type = string
}
variable "google_client_secret" {
  type = string
}

variable "image_name" {
  type = string
}

variable "project_id" {
  type = string
}


provider "google" {
  project = var.project_id
  region  = "us-west1"
}

resource "google_project_service" "run" {
  service = "run.googleapis.com"
}

resource "google_cloud_run_service" "my-service" {
  name     = "gourd-backend"
  location = "us-west1"

  template {
    spec {
      containers {
        resources {
          limits = {
            memory = "1G"
          }
        }
        image = var.image_name
        env {
          name  = "DB_DBNAME"
          value = "food"
        }
        env {
          name  = "DB_HOST"
          value = "34.66.204.3"
        }
        env {
          name  = "DB_PORT"
          value = "5432"
        }
        env {
          name  = "DB_USER"
          value = "postgres"
        }
        env {
          name  = "DB_PASSWORD"
          value = var.db_password
        }
        env {
          name  = "GOOGLE_CLOUD_PROJECT"
          value = var.project_id
        }
        env {
          name  = "JWT_KEY"
          value = "tmp"
        }
        env {
          name  = "GOOGLE_CLIENT_ID"
          value = var.google_client_id
        }
        env {
          name  = "GOOGLE_CLIENT_SECRET"
          value = var.google_client_secret
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [google_project_service.run]
}

resource "google_cloud_run_service_iam_member" "allUsers" {
  service  = google_cloud_run_service.my-service.name
  location = google_cloud_run_service.my-service.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

resource "google_service_account" "sa" {
  account_id   = "my-service-account"
  display_name = "A service account that only Jane can use"
}

data "google_iam_policy" "admin" {
  binding {
    role = "roles/run.invoker"

    members = [
      "serviceAccount:${google_service_account.sa.email}",
    ]
  }
}
data "google_iam_policy" "admin2" {
  binding {
    role = "roles/cloudscheduler.admin"

    members = [
      "serviceAccount:${google_service_account.sa.email}",
    ]
  }
}

# https://benjamincongdon.me/blog/2019/11/21/Setting-up-Cloud-Scheduler-to-Trigger-Cloud-Run/

# resource "google_service_account_iam_binding" "admin-account-iam" {
#   service_account_id = google_service_account.sa.name
#   role               = "roles/run.invoker"

#   members = [
#     "user:14nicholasse@gmail.com",
#   ]
# }





output "url" {
  value = google_cloud_run_service.my-service.status[0].url
}

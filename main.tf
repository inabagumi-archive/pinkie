terraform {
  cloud {
    organization = "inabagumi"

    workspaces {
      name = "pinkie"
    }
  }

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.39.0"
    }

    google-beta = {
      source  = "hashicorp/google-beta"
      version = "4.38.0"
    }

    github = {
      source  = "integrations/github"
      version = "5.3.0"
    }
  }
}

provider "google" {
  project = var.project
  region  = var.region
}

provider "google-beta" {
  project = var.project
  region  = var.region
}

provider "github" {
  owner = var.repo_owner
}

locals {
  fetch_job_name = "fetch"
}

data "google_project" "project" {}

resource "google_service_account" "terraform" {
  account_id   = "terraform"
  description  = "Serivce Account for Terraform"
  display_name = "Terraform"
  project = var.project
}

resource "google_service_account" "gha" {
  account_id   = "github-actions"
  description  = "Service Account for GitHub Actions"
  display_name = "GitHub Actions"
  project = var.project
}

resource "google_service_account" "pinkie" {
  account_id   = "pinkie"
  description  = "Service Account for Pinkie"
  display_name = "Pinkie"
  project = var.project
}

resource "google_project_iam_binding" "artifactregistry_admin" {
  members = ["serviceAccount:${google_service_account.terraform.email}"]
  project = var.project
  role    = "roles/artifactregistry.admin"
}

resource "google_project_iam_binding" "cloudscheduler_admin" {
  members = ["serviceAccount:${google_service_account.terraform.email}"]
  project = var.project
  role    = "roles/cloudscheduler.admin"
}

resource "google_project_iam_binding" "resourcemanager_project_iam_admin" {
  members = ["serviceAccount:${google_service_account.terraform.email}"]
  project = var.project
  role    = "roles/resourcemanager.projectIamAdmin"
}

resource "google_project_iam_binding" "iam_service_account_admin" {
  members = ["serviceAccount:${google_service_account.terraform.email}"]
  project = var.project
  role    = "roles/iam.serviceAccountAdmin"
}

resource "google_project_iam_binding" "iam_service_account_user" {
  members = [
    "serviceAccount:${google_service_account.terraform.email}",
    "serviceAccount:${google_service_account.gha.email}",
    "serviceAccount:${google_service_account.pinkie.email}",
  ]
  project = var.project
  role    = "roles/iam.serviceAccountUser"
}

resource "google_project_iam_binding" "iam_workload_identity_pool_admin" {
  members = [
    "serviceAccount:${google_service_account.terraform.email}",
    "serviceAccount:${google_service_account.gha.email}",
  ]
  project = var.project
  role    = "roles/iam.workloadIdentityPoolAdmin"
}

module "gh_oidc" {
  source  = "terraform-google-modules/github-actions-runners/google//modules/gh-oidc"
  version = "3.1.0"

  pool_id               = "${var.repo_name}-pool"
  project_id            = var.project
  provider_display_name = "GitHub"
  provider_id           = "github"
  sa_mapping = {
    (google_service_account.gha.account_id) = {
      attribute = "attribute.repository/${var.repo_owner}/${var.repo_name}"
      sa_name   = google_service_account.gha.name
    }
  }
}

resource "google_artifact_registry_repository" "containers" {
  provider = google-beta

  format        = "DOCKER"
  location      = var.region
  project       = var.project
  repository_id = "containers"
}

resource "google_artifact_registry_repository_iam_member" "gha" {
  provider = google-beta

  location   = google_artifact_registry_repository.containers.location
  member     = "serviceAccount:${google_service_account.gha.email}"
  project    = google_artifact_registry_repository.containers.project
  repository = google_artifact_registry_repository.containers.name
  role       = "roles/artifactregistry.writer"
}

resource "google_artifact_registry_repository_iam_member" "reader" {
  provider = google-beta

  location   = google_artifact_registry_repository.containers.location
  member     = "allUsers"
  project    = google_artifact_registry_repository.containers.project
  repository = google_artifact_registry_repository.containers.name
  role       = "roles/artifactregistry.reader"
}

resource "google_cloud_scheduler_job" "fetch" {
  attempt_deadline = "600s"
  name             = local.fetch_job_name
  project          = var.project
  region           = var.region
  schedule         = "1-51/10 * * * *"
  time_zone        = "UTC"

  http_target {
    http_method = "POST"
    uri         = "https://${var.region}-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/${var.project}/jobs/${local.fetch_job_name}:run"

    oauth_token {
      service_account_email = google_service_account.pinkie.email
    }
  }

  retry_config {
    retry_count = 0
  }
}

resource "github_actions_secret" "project" {
  plaintext_value = var.project
  repository      = var.repo_name
  secret_name     = "GOOGLE_PROJECT"
}

resource "github_actions_secret" "region" {
  plaintext_value = var.region
  repository      = var.repo_name
  secret_name     = "GOOGLE_REGION"
}

resource "github_actions_secret" "service_account" {
  plaintext_value = google_service_account.gha.email
  repository      = var.repo_name
  secret_name     = "GOOGLE_SERVICE_ACCOUNT"
}

resource "github_actions_secret" "workload_identity_provider" {
  plaintext_value = module.gh_oidc.provider_name
  repository      = var.repo_name
  secret_name     = "GOOGLE_WORKLOAD_IDENTITY_PROVIDER"
}

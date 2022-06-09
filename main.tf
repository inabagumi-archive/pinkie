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
      version = "4.24.0"
    }

    google-beta = {
      source  = "hashicorp/google-beta"
      version = "4.24.0"
    }
  }
}

provider "google" {
  project = var.project
}

provider "google-beta" {
  project = var.project
}

provider "github" {
  owner = var.repo_owner
}

resource "google_service_account" "gha" {
  account_id   = "github-actions"
  display_name = "Service Account for GitHub Actions"
}

resource "google_project_iam_binding" "iam_workload_identity_pool_admin" {
  members = ["serviceAccount:${google_service_account.gha.email}"]
  project = var.project
  role    = "roles/iam.workloadIdentityPoolAdmin"
}

resource "google_project_iam_binding" "iam_service_account_admin" {
  members = ["serviceAccount:${google_service_account.gha.email}"]
  project = var.project
  role    = "roles/iam.serviceAccountAdmin"
}

module "gh_oidc" {
  source  = "terraform-google-modules/github-actions-runners/google//modules/gh-oidc"
  version = "3.0.0"

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

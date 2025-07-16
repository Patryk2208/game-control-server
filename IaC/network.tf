resource "google_compute_network" "game_vpc" {
  name                    = "game-vpc"
  auto_create_subnetworks = false
  description             = "VPC for game platform components"
}

resource "google_compute_subnetwork" "game_subnet" {
  name          = "game-subnet"
  ip_cidr_range = "10.10.1.0/24"
  region        = var.region
  network       = google_compute_network.game_vpc.id

  secondary_ip_range {
    range_name    = "pods"
    ip_cidr_range = "10.20.0.0/16"
  }

  secondary_ip_range {
    range_name    = "services"
    ip_cidr_range = "10.30.0.0/20"
  }
}


resource "google_compute_router" "nat_router" {
  name    = "game-nat-router"
  network = google_compute_network.game_vpc.id
  region  = var.region
}


resource "google_compute_router_nat" "game_nat" {
  name                               = "game-platform-nat"
  router                             = google_compute_router.nat_router.name
  region                             = var.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"

  log_config {
    enable = true
    filter = "ERRORS_ONLY"
  }
}

# FIREWALL
resource "google_compute_firewall" "allow_internal" {
  name        = "allow-internal"
  network     = google_compute_network.game_vpc.name
  direction   = "INGRESS"

  allow {
    protocol = "all"
  }

  source_ranges = [
    google_compute_subnetwork.game_subnet.ip_cidr_range,
    google_compute_subnetwork.game_subnet.secondary_ip_range[0].ip_cidr_range,
    google_compute_subnetwork.game_subnet.secondary_ip_range[1].ip_cidr_range
  ]
}


resource "google_compute_firewall" "allow_websocket" {
  name        = "allow-websocket"
  network     = google_compute_network.game_vpc.name
  direction   = "INGRESS"
  priority    = 1000

  allow {
    protocol = "tcp"
    ports    = ["80", "443"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["matchmaking"]
}

resource "google_compute_firewall" "allow_nodeports" {
  name        = "allow-nodeports"
  network     = google_compute_network.game_vpc.name
  direction   = "INGRESS"
  priority    = 1000

  allow {
    protocol = "tcp"
    ports    = ["30000-32767"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["test-node"]
}


resource "google_compute_firewall" "allow_game_tcp" {
  name        = "allow-game-tcp"
  network     = google_compute_network.game_vpc.name
  direction   = "INGRESS"
  priority    = 1000

  allow {
    protocol = "tcp"
    ports    = ["7000-8000"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["game-server"]
}


/*resource "google_compute_firewall" "allow_ssh" {
  name        = "allow-ssh"
  network     = google_compute_network.game_vpc.name
  direction   = "INGRESS"

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["/32"]  //todo
  target_tags   = ["ssh"]
}*/
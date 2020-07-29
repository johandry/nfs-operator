resource "ibm_resource_group" "group" {
  name = "${var.project_name}-${var.environment}-group"
}

// data "ibm_resource_group" "group" {
//   name = "${var.project_name}-${var.environment}-group"
// }

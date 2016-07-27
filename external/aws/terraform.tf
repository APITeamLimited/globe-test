variable "access_key" ***REMOVED******REMOVED***
variable "secret_key" ***REMOVED******REMOVED***
variable "region" ***REMOVED*** default = "eu-west-1" ***REMOVED***
variable "ami" ***REMOVED*** default = "ami-a4d44ed7" ***REMOVED***
variable "key_name" ***REMOVED*** default = "speedboat-test" ***REMOVED***

output "loadgen_ip" ***REMOVED***
	value = "$***REMOVED***aws_instance.loadgen.public_ip***REMOVED***"
***REMOVED***
output "influx_ip" ***REMOVED***
	value = "$***REMOVED***aws_instance.influx.public_ip***REMOVED***"
***REMOVED***
output "web_ip" ***REMOVED***
	value = "$***REMOVED***aws_instance.web.public_ip***REMOVED***"
***REMOVED***

provider "aws" ***REMOVED***
	access_key = "$***REMOVED***var.access_key***REMOVED***"
	secret_key = "$***REMOVED***var.secret_key***REMOVED***"
	region = "$***REMOVED***var.region***REMOVED***"
***REMOVED***

resource "aws_security_group" "group" ***REMOVED***
	name = "speedboat-test"
	description = "Security group for Speedboat test setups"
	
	ingress ***REMOVED***
		from_port = 0
		to_port = 0
		protocol = "-1"
		cidr_blocks = ["0.0.0.0/0"]
	***REMOVED***
	
	egress ***REMOVED***
		from_port = 0
		to_port = 0
		protocol = "-1"
		cidr_blocks = ["0.0.0.0/0"]
	***REMOVED***
***REMOVED***

resource "aws_placement_group" "group" ***REMOVED***
	name = "speedboat-test"
	strategy = "cluster"
***REMOVED***

resource "aws_instance" "loadgen" ***REMOVED***
	instance_type = "m4.xlarge"
	ami = "$***REMOVED***var.ami***REMOVED***"
	placement_group = "$***REMOVED***aws_placement_group.group.name***REMOVED***"
	security_groups = ["$***REMOVED***aws_security_group.group.name***REMOVED***"]
	key_name = "$***REMOVED***var.key_name***REMOVED***"
	tags ***REMOVED***
		Name = "sbt-loadgen"
	***REMOVED***
	ebs_optimized = true

	connection ***REMOVED***
		user = "ubuntu"
		private_key = "$***REMOVED***file("$***REMOVED***var.key_name***REMOVED***.pem")***REMOVED***"
	***REMOVED***
	provisioner "remote-exec" ***REMOVED***
		inline = [
			"mkdir -p /home/ubuntu/go/src/github.com/loadimpact/speedboat",
			"echo 'export GOPATH=$HOME/go' >> /home/ubuntu/.profile",
			"echo 'export PATH=$PATH:$GOPATH/bin' >> /home/ubuntu/.profile",
			"sudo mkdir -p /etc/salt",
			"sudo ln -s /home/ubuntu/go/src/github.com/loadimpact/speedboat/external/aws/salt/master.yml /etc/salt/master",
			"sudo ln -s /home/ubuntu/go/src/github.com/loadimpact/speedboat/external/aws/salt/grains_loadgen.yml /etc/salt/grains",
		]
	***REMOVED***
	provisioner "file" ***REMOVED***
		source = "../../"
		destination = "/home/ubuntu/go/src/github.com/loadimpact/speedboat"
	***REMOVED***
	provisioner "remote-exec" ***REMOVED***
		inline = [
			"curl -L https://bootstrap.saltstack.com | sudo sh -s -- -n -M -A 127.0.0.1 -i loadgen stable 2016.3.1",
		]
	***REMOVED***
***REMOVED***

resource "aws_instance" "influx" ***REMOVED***
	instance_type = "m4.xlarge"
	ami = "$***REMOVED***var.ami***REMOVED***"
	placement_group = "$***REMOVED***aws_placement_group.group.name***REMOVED***"
	security_groups = ["$***REMOVED***aws_security_group.group.name***REMOVED***"]
	key_name = "$***REMOVED***var.key_name***REMOVED***"
	tags ***REMOVED***
		Name = "sbt-influx"
	***REMOVED***
	ebs_optimized = true

	connection ***REMOVED***
		user = "ubuntu"
		private_key = "$***REMOVED***file("$***REMOVED***var.key_name***REMOVED***.pem")***REMOVED***"
	***REMOVED***
	provisioner "remote-exec" ***REMOVED***
		inline = [
			"sudo mkdir -p /etc/salt",
			"sudo touch /etc/salt/grains",
			"sudo chown ubuntu:ubuntu /etc/salt/grains",
		]
	***REMOVED***
	provisioner "file" ***REMOVED***
		source = "salt/grains_influx.yml"
		destination = "/etc/salt/grains"
	***REMOVED***
	provisioner "remote-exec" ***REMOVED***
		inline = [
			"curl -L https://bootstrap.saltstack.com | sudo sh -s -- -n -A $***REMOVED***aws_instance.loadgen.private_ip***REMOVED*** -i influx stable 2016.3.1"
		]
	***REMOVED***
***REMOVED***

resource "aws_instance" "web" ***REMOVED***
	instance_type = "m4.xlarge"
	ami = "$***REMOVED***var.ami***REMOVED***"
	placement_group = "$***REMOVED***aws_placement_group.group.name***REMOVED***"
	security_groups = ["$***REMOVED***aws_security_group.group.name***REMOVED***"]
	key_name = "$***REMOVED***var.key_name***REMOVED***"
	tags ***REMOVED***
		Name = "sbt-web"
	***REMOVED***
	ebs_optimized = true

	connection ***REMOVED***
		user = "ubuntu"
		private_key = "$***REMOVED***file("$***REMOVED***var.key_name***REMOVED***.pem")***REMOVED***"
	***REMOVED***
	provisioner "remote-exec" ***REMOVED***
		inline = [
			"sudo mkdir -p /etc/salt",
			"sudo touch /etc/salt/grains",
			"sudo chown ubuntu:ubuntu /etc/salt/grains",
		]
	***REMOVED***
	provisioner "file" ***REMOVED***
		source = "salt/grains_web.yml"
		destination = "/etc/salt/grains"
	***REMOVED***
	provisioner "remote-exec" ***REMOVED***
		inline = [
			"curl -L https://bootstrap.saltstack.com | sudo sh -s -- -n -A $***REMOVED***aws_instance.loadgen.private_ip***REMOVED*** -i web stable 2016.3.1"
		]
	***REMOVED***
***REMOVED***

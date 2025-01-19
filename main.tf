provider "aws" {
  region = "us-east-1"
}

resource "aws_instance" "example" {
  ami           = "ami-0c55b159cbfafe1f0" # Replace with a valid AMI ID
  instance_type = "t2.micro"

  tags = {
    Name = "ExampleInstance"
  }
}

resource "aws_ec2_host" "name" {
  availability_zone = "us-east-1"
}

resource "aws_s3_bucket" "create-bucket" {
  bucket = "my-bucket"
  tags   = {
    Name        = "my-bucket"
    Environment = "dev"
  }
}

resource "aws_s3_bucket" "test-bucket" {
  bucket = "test-bucket"

  tags = {
    Name = "test-bucket"
    Environment = "dev"
  }
}
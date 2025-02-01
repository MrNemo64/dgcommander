# Auto-generated Makefile. Do not edit manually.

generate.makefile-debug:
	cd tools/makefile_debug && go run .

# If we kill a debug process from VSCode the commands are not deleted. This is to clean up for when this happens.
delete.commands:
	cd tools/delete_commands && go run .

example.message:
	cd examples/message && go run .

example.middleware:
	cd examples/middleware && go run .

example.slash-calculator:
	cd examples/slash-calculator && go run .

example.tasks:
	cd examples/tasks && go run .

example.users:
	cd examples/users && go run .


common.queueb.org/tools
=======================

.. contents::
  :local:
  :depth: 2

Tools mostly contains application level tooling which might be
considered as generic.

The package provides simple primitives to work with:

- env, env_helpers -- cli-related helpers to work with application flags.
  using environment variables. See examples/main.go for details.
- maps -- maps related helpers.
- os -- operating system related helpers.
- tools -- code helpers (experimental).

It's still in development and used for internal queueb projects.

Installation
------------

Install the package with:

.. code-block:: console

   go get common.queueb.org/tools

Usage
-----

Import the package:

.. code-block:: go

   import "common.queueb.org/tools"

Purpose
-------
Application level primitives.

Development
-----------

Documentation
~~~~~~~~~~~~~
Documentation is in development process. Sooner or later it would be available on 
Go Package Documentation: https://pkg.go.dev/common.queueb.org/tools

Source code: https://github.com/queueb-org/common-tools 

Testing
~~~~~~~

.. code-block:: bash

  $ go test -coverprofile=.coverage ./... && go tool cover -func=.coverage

License
-------

See the ``LICENSE`` file in the repository.
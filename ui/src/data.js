var data = [
    {
        "ID": 0,
        "CreatedAt": "0001-01-01T00:00:00Z",
        "UpdatedAt": "0001-01-01T00:00:00Z",
        "DeletedAt": null,
        "Server": "https://github.com",
        "Namespace": "geerlingguy",
        "Repo": "ansible-role-jenkins",
        "Meta": {
            "Dependencies": [
                "geerlingguy.java"
            ],
            "GalaxyInfo": {
                "Author": "geerlingguy",
                "Description": "Jenkins CI",
                "Company": "Midwestern Mac, LLC",
                "License": "license (BSD, MIT)",
                "MinAnsibleVersion": 2.4,
                "Platforms": [
                    {
                        "Name": "EL",
                        "Versions": [
                            "6",
                            "7"
                        ]
                    },
                    {
                        "Name": "Fedora",
                        "Versions": [
                            "all"
                        ]
                    },
                    {
                        "Name": "Debian",
                        "Versions": [
                            "all"
                        ]
                    },
                    {
                        "Name": "Ubuntu",
                        "Versions": [
                            "all"
                        ]
                    }
                ],
                "GalaxyTags": [
                    "development",
                    "packaging",
                    "jenkins",
                    "ci"
                ]
            }
        },
        "Readme": "# Ansible Role: Jenkins CI\n\n[![Build Status](https://travis-ci.org/geerlingguy/ansible-role-jenkins.svg?branch=master)](https://travis-ci.org/geerlingguy/ansible-role-jenkins)\n\nInstalls Jenkins CI on RHEL/CentOS and Debian/Ubuntu servers.\n\n## Requirements\n\nRequires `curl` to be installed on the server. Also, newer versions of Jenkins require Java 8+ (see the test playbooks inside the `molecule/default` directory for an example of how to use newer versions of Java for your OS).\n\n## Role Variables\n\nAvailable variables are listed below, along with default values (see `defaults/main.yml`):\n\n    jenkins_package_state: present\n\nThe state of the `jenkins` package install. By default this role installs Jenkins but will not upgrade Jenkins (when using package-based installs). If you want to always update to the latest version, change this to `latest`.\n\n    jenkins_hostname: localhost\n\nThe system hostname; usually `localhost` works fine. This will be used during setup to communicate with the running Jenkins instance via HTTP requests.\n\n    jenkins_home: /var/lib/jenkins\n\nThe Jenkins home directory which, amongst others, is being used for storing artifacts, workspaces and plugins. This variable allows you to override the default `/var/lib/jenkins` location.\n\n    jenkins_http_port: 8080\n\nThe HTTP port for Jenkins' web interface.\n\n    jenkins_admin_username: admin\n    jenkins_admin_password: admin\n\nDefault admin account credentials which will be created the first time Jenkins is installed.\n\n    jenkins_admin_password_file: \"\"\n\nDefault admin password file which will be created the first time Jenkins is installed as /var/lib/jenkins/secrets/initialAdminPassword\n\n    jenkins_jar_location: /opt/jenkins-cli.jar\n\nThe location at which the `jenkins-cli.jar` jarfile will be kept. This is used for communicating with Jenkins via the CLI.\n\n    jenkins_plugins: []\n\nJenkins plugins to be installed automatically during provisioning.\n\n    jenkins_plugins_install_dependencies: true\n\nWhether Jenkins plugins to be installed should also install any plugin dependencies.\n\n    jenkins_plugins_state: present\n\nUse `latest` to ensure all plugins are running the most up-to-date version.\n\n    jenkins_plugin_updates_expiration: 86400\n\nNumber of seconds after which a new copy of the update-center.json file is downloaded. Set it to 0 if no cache file should be used.\n\n    jenkins_updates_url: \"https://updates.jenkins.io\"\n\nThe URL to use for Jenkins plugin updates and update-center information.\n\n    jenkins_plugin_timeout: 30\n\nThe server connection timeout, in seconds, when installing Jenkins plugins.\n\n    jenkins_version: \"1.644\"\n    jenkins_pkg_url: \"http://www.example.com\"\n\n(Optional) Then Jenkins version can be pinned to any version available on `http://pkg.jenkins-ci.org/debian/` (Debian/Ubuntu) or `http://pkg.jenkins-ci.org/redhat/` (RHEL/CentOS). If the Jenkins version you need is not available in the default package URLs, you can override the URL with your own; set `jenkins_pkg_url` (_Note_: the role depends on the same naming convention that `http://pkg.jenkins-ci.org/` uses).\n\n    jenkins_url_prefix: \"\"\n\nUsed for setting a URL prefix for your Jenkins installation. The option is added as `--prefix={{ jenkins_url_prefix }}` to the Jenkins initialization `java` invocation, so you can access the installation at a path like `http://www.example.com{{ jenkins_url_prefix }}`. Make sure you start the prefix with a `/` (e.g. `/jenkins`).\n\n    jenkins_connection_delay: 5\n    jenkins_connection_retries: 60\n\nAmount of time and number of times to wait when connecting to Jenkins after initial startup, to verify that Jenkins is running. Total time to wait = `delay` * `retries`, so by default this role will wait up to 300 seconds before timing out.\n\n    # For RedHat/CentOS (role default):\n    jenkins_repo_url: http://pkg.jenkins-ci.org/redhat/jenkins.repo\n    jenkins_repo_key_url: http://pkg.jenkins-ci.org/redhat/jenkins-ci.org.key\n    # For Debian (role default):\n    jenkins_repo_url: deb http://pkg.jenkins-ci.org/debian binary/\n    jenkins_repo_key_url: http://pkg.jenkins-ci.org/debian/jenkins-ci.org.key\n\nThis role will install the latest version of Jenkins by default (using the official repositories as listed above). You can override these variables (use the correct set for your platform) to install the current LTS version instead:\n\n    # For RedHat/CentOS LTS:\n    jenkins_repo_url: http://pkg.jenkins-ci.org/redhat-stable/jenkins.repo\n    jenkins_repo_key_url: http://pkg.jenkins-ci.org/redhat-stable/jenkins-ci.org.key\n    # For Debian/Ubuntu LTS:\n    jenkins_repo_url: deb http://pkg.jenkins-ci.org/debian-stable binary/\n    jenkins_repo_key_url: http://pkg.jenkins-ci.org/debian-stable/jenkins-ci.org.key\n\nIt is also possible stop the repo file being added by setting  `jenkins_repo_url = ''`. This is useful if, for example, you sign your own packages or run internal package management (e.g. Spacewalk).\n\n    jenkins_java_options: \"-Djenkins.install.runSetupWizard=false\"\n\nExtra Java options for the Jenkins launch command configured in the init file can be set with the var `jenkins_java_options`. For example, if you want to configure the timezone Jenkins uses, add `-Dorg.apache.commons.jelly.tags.fmt.timeZone=America/New_York`. By default, the option to disable the Jenkins 2.0 setup wizard is added.\n\n    jenkins_init_changes:\n      - option: \"JENKINS_ARGS\"\n        value: \"--prefix={{ jenkins_url_prefix }}\"\n      - option: \"JENKINS_JAVA_OPTIONS\"\n        value: \"{{ jenkins_java_options }}\"\n\nChanges made to the Jenkins init script; the default set of changes set the configured URL prefix and add in configured Java options for Jenkins' startup. You can add other option/value pairs if you need to set other options for the Jenkins init file.\n\n## Dependencies\n\n  - geerlingguy.java\n\n## Example Playbook\n\n```yaml\n- hosts: jenkins\n  vars:\n    jenkins_hostname: jenkins.example.com\n  roles:\n    - role: geerlingguy.java\n      become: yes\n    - role: geerlingguy.jenkins\n      become: yes\n```\n\n## License\n\nMIT (Expat) / BSD\n\n## Author Information\n\nThis role was created in 2014 by [Jeff Geerling](https://www.jeffgeerling.com/), author of [Ansible for DevOps](https://www.ansiblefordevops.com/).\n",
        "Supported": false,
        "Rated": 0
    },
    {
        "ID": 0,
        "CreatedAt": "0001-01-01T00:00:00Z",
        "UpdatedAt": "0001-01-01T00:00:00Z",
        "DeletedAt": null,
        "Server": "https://github.com",
        "Namespace": "geerlingguy",
        "Repo": "ansible-role-java",
        "Meta": {
            "Dependencies": [],
            "GalaxyInfo": {
                "Author": "geerlingguy",
                "Description": "Java for Linux",
                "Company": "Midwestern Mac, LLC",
                "License": "license (BSD, MIT)",
                "MinAnsibleVersion": 2.4,
                "Platforms": [
                    {
                        "Name": "EL",
                        "Versions": [
                            "6",
                            "7"
                        ]
                    },
                    {
                        "Name": "Fedora",
                        "Versions": [
                            "all"
                        ]
                    },
                    {
                        "Name": "Debian",
                        "Versions": [
                            "wheezy",
                            "jessie",
                            "stretch"
                        ]
                    },
                    {
                        "Name": "Ubuntu",
                        "Versions": [
                            "precise",
                            "trusty",
                            "xenial",
                            "bionic"
                        ]
                    },
                    {
                        "Name": "FreeBSD",
                        "Versions": [
                            "10.2"
                        ]
                    }
                ],
                "GalaxyTags": [
                    "development",
                    "system",
                    "web",
                    "java",
                    "jdk",
                    "openjdk",
                    "oracle"
                ]
            }
        },
        "Readme": "# Ansible Role: Java\n\n[![Build Status](https://travis-ci.org/geerlingguy/ansible-role-java.svg?branch=master)](https://travis-ci.org/geerlingguy/ansible-role-java)\n\nInstalls Java for RedHat/CentOS and Debian/Ubuntu linux servers.\n\n## Requirements\n\nNone.\n\n## Role Variables\n\nAvailable variables are listed below, along with default values:\n\n    # The defaults provided by this role are specific to each distribution.\n    java_packages:\n      - java-1.7.0-openjdk\n\nSet the version/development kit of Java to install, along with any other necessary Java packages. Some other options include are included in the distribution-specific files in this role's 'defaults' folder.\n\n    java_home: \"\"\n\nIf set, the role will set the global environment variable `JAVA_HOME` to this value.\n\n## Dependencies\n\nNone.\n\n## Example Playbook (using default package, usually OpenJDK 7)\n\n    - hosts: servers\n      roles:\n        - role: geerlingguy.java\n          become: yes\n\n## Example Playbook (install OpenJDK 8)\n\nFor RHEL / CentOS:\n\n    - hosts: server\n      roles:\n        - role: geerlingguy.java\n          when: \"ansible_os_family == 'RedHat'\"\n          java_packages:\n            - java-1.8.0-openjdk\n\nFor Ubuntu < 16.04:\n\n    - hosts: server\n      tasks:\n        - name: installing repo for Java 8 in Ubuntu\n  \t      apt_repository: repo='ppa:openjdk-r/ppa'\n    \n    - hosts: server\n      roles:\n        - role: geerlingguy.java\n          when: \"ansible_os_family == 'Debian'\"\n          java_packages:\n            - openjdk-8-jdk\n\n## License\n\nMIT / BSD\n\n## Author Information\n\nThis role was created in 2014 by [Jeff Geerling](https://www.jeffgeerling.com/), author of [Ansible for DevOps](https://www.ansiblefordevops.com/).\n",
        "Supported": false,
        "Rated": 0
    }
]


const userData = {
    firstname:"Robert",
    lastname:"Fletcher",
    username:"username"
}


export let roles = data;
export let user = userData
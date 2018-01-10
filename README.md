# Luzifer / ansible-role-version

Very simple tool to update a `requirements.yml` file for [Ansible Galaxy](https://galaxy.ansible.com/) with specific versions of roles.

## Usage

Lets say you do have a repository containing this `requirements.yml`:

```yaml
---

- name: docker
  src: https://github.com/luzifer-ansible/docker
  version: v0.1.0
- name: docker-compose
  src: https://github.com/luzifer-ansible/docker-compose
  version: v1.0.0

...
```

Now your CI system should update the version of `docker-compose` to `v1.0.1` and you don't want to fiddle with bash magic:

```bash
$ ansible-role-version set docker-compose v1.0.1
```

And you're done!

Sure, this example is a bit constructed: In reality this tool was written to update a bunch of different repositories each having a way bigger list of roles to include and the tool is used in a script to update all of them and create pull-requests out of the change.

# Redeploy applications within a namespace on a schedule

> This is an experiment and Okteto does not officially support it.

- Create an [Okteto Admin Token](https://www.okteto.com/docs/admin/dashboard/#admin-access-tokens)

- Export the token to a local variable:

```bash
export OKTETO_ADMIN_TOKEN=<<your-token>>
```

- Create a namespace, and, via the admin section, mark it as [Keep awake](https://www.okteto.com/docs/admin/dashboard/#namespaces)

- Export the namespace name to a local variable:

```bash
export NAMESPACE=<<your-namespace>>
```

- Create a local variable to define the redeploy cronjob schedule:

```bash
export REDEPLOY_JOB_SCHEDULE="0 20 * * *"
```

For example, 0 0 13 * 5 states that the task must be started every Friday at midnight, as well as on the 13th of each month at midnight.

- Create another local variable to define the repository to redeploy:

```bash
export TARGET_REPOSITORY=https://github.com/okteto/movies
```

Optionally, you can specify also the branch to filter defining the env var `TARGET_BRANCH`.

```bash
export TARGET_BRANCH=main
```

- Run the following command to create the cronjob:

```bash
okteto deploy -n ${NAMESPACE} --var OKTETO_ADMIN_TOKEN=${OKTETO_ADMIN_TOKEN} --var REDEPLOY_JOB_SCHEDULE=${REDEPLOY_JOB_SCHEDULE} --var TARGET_REPOSITORY=${TARGET_REPOSITORY} --var TARGET_BRANCH=${TARGET_BRANC}
```

## Force the execution of the job

To force the execution of the redeploy namespaces job, run the following commands:

```bash
okteto kubeconfig
kubectl -n ${NAMESPACE} create job --from=cronjob/redeploy-apps redeploy-apps-$(date +%s)
```
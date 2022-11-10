<h2>
Globe Test
</h2>

##### Run distributed load tests using the K6 runtime engine

Actively being developed by APITeam (<a href="https://apiteam.cloud">https://apiteam.cloud</a>). APITeam is an all in one platform for designing, testing and scaling your APIs collaboratively. 

Note: GlobeTest and APITeam are not affiliated with the K6 project.

The aim of this project is to provide a simple way to run distributed load tests using the K6 runtime engine that avoids the need for Kubernetes. A further aim is to provide simultanous concurrent execution of multiple load tests.

Globe Test is designed to be deployed as several different containers on different hosts. The two kinds of nodes are orchestrator and worker nodes. The orchestrator nodes are responsible for managing the load tests and the worker nodes are responsible for executing the load tests.

Communciation between nodes is accomplished via several redis containers, logs and metrics are stored in a gridfs database.

More documentation on how to use this shall be added soon.
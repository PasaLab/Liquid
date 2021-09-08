<!-- msg modal -->
<div class="modal fade" id="modal-msg" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
	<div class="modal-dialog">
		<div class="modal-content panel-warning">
			<div class="modal-header panel-heading">
				<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span>
				</button>
				<h4 id="modal-msg-title" class="modal-title">Notice</h4>
			</div>
			<div class="modal-body">
				<h4 id="modal-msg-content" class="text-msg text-center">Something is wrong!</h4>
			</div>
		</div>
	</div>
</div>

<!-- job description modal -->
<div class="modal fade" id="modal-job-description" tabindex="-1" role="dialog" aria-labelledby="myModalLabel"
     aria-hidden="true">
	<div class="modal-dialog">
		<div class="modal-content panel-info">
			<div class="modal-header panel-heading">
				<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">Describe This Job</h4>
			</div>
			<div class="modal-body">
				<pre id="modal-job-description-content"></pre>
			</div>
		</div>
	</div>
</div>

<!-- node GPU detail modal -->
<div class="modal fade" id="modal-resource-gpu-detail" tabindex="-1" role="dialog" aria-labelledby="myModalLabel"
     aria-hidden="true">
	<div class="modal-dialog">
		<div class="modal-content panel-info">
			<div class="modal-header panel-heading">
				<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">GPUs on this node</h4>
			</div>
			<div class="modal-body">
				<pre id="modal-resource-gpu-detail-content"></pre>
			</div>
		</div>
	</div>
</div>

<!-- task logs modal -->
<div class="modal fade" id="modal-task-logs" tabindex="-1" role="dialog" aria-labelledby="myModalLabel"
     aria-hidden="true">
	<div class="modal-dialog modal-lg">
		<div class="modal-content panel-info">
			<div class="modal-header panel-heading">
				<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">Task Outputs</h4>
			</div>
			<div class="modal-body">
				<pre id="modal-task-logs-content"></pre>
			</div>
		</div>
	</div>
</div>

<!-- job modal -->
<div class="modal fade" id="modal-job" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
	<div class="modal-dialog modal-lg">
		<div class="modal-content panel-info">
			<div class="modal-header panel-heading">
				<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span>
				</button>
				<h4 id="modal-job-title" class="modal-title">Submit New Job</h4>
			</div>
			<div class="modal-body">
				<form class="form" action="javascript:void(0)">
					<label>Job Name</label>
					<div class="form-group form-group-lg">
						<label for="form-job-name" class="sr-only">Job Name</label>
						<input type="text" id="form-job-name" class="form-control" maxlength="64"
						       placeholder="A readable job name" required/>
					</div>
					<label>Input Dir</label>
					<div class="form-group form-group-lg">
						<label for="form-job-workspace" class="sr-only">Workspace</label>
						<select id="form-job-workspace" class="form-control">
							<option value="">None</option>
						</select>
					</div>
					<label>Model Dir</label>
					<div class="form-group form-group-lg">
						<label for="form-job-model-dir" class="sr-only">Model Dir</label>
						<input type="text" id="form-job-model-dir" class="form-control" maxlength="256"
						       placeholder="Dir for model checkpoints" required/>
					</div>
					<label>Output Dir</label>
					<div class="form-group form-group-lg">
						<label for="form-job-output-dir" class="sr-only">Workspace</label>
						<input type="text" id="form-job-output-dir" class="form-control" maxlength="256"
						       placeholder="Dir for result data" required/>
					</div>
					<label>Queue</label>
					<div class="form-group form-group-lg">
						<label for="form-job-cluster" class="sr-only">Virtual Cluster</label>
						<select id="form-job-cluster" class="form-control">
							<option value="1">default</option>
						</select>
					</div>

					<label>Priority</label>
					<div class="form-group form-group-lg">
						<label for="form-job-priority" class="sr-only">Job Priority</label>
						<select id="form-job-priority" class="form-control">
							<option value="99">Urgent</option>
							<option value="50">High</option>
							<option value="25" selected>Medium</option>
							<option value="1">Low</option>
						</select>
					</div>
					<label class="hidden">Locality</label>
					<div class="form-group form-group-lg hidden">
						<label for="form-job-locality" class="sr-only">Locality</label>
						<select id="form-job-locality" class="form-control">
							<option value="1">Positive</option>
							<option value="0" selected>Any</option>
							<option value="-1">Negative</option>
						</select>
					</div>
					<label class="hidden">Run Before</label>
					<div class="form-group form-group-lg hidden">
						<div class='input-group date date-picker'>
							<label for="form-job-run-before" class="sr-only">Run Before</label>
							<input type='text' class="form-control" placeholder="Run this job before"
							       id="form-job-run-before"
							       autocomplete="off"/>
							<div class="input-group-addon">
								<span class="glyphicon glyphicon-calendar"></span>
							</div>
						</div>
					</div>
					<label>Tasks</label>
					<div id="form-job-tasks">
						<div class="row">
							<div class="col-md-4">
								<label>Docker Image</label>
								<input type="text" class="form-control task-image" maxlength=""
								       value="quickdeploy/yao-tensorflow:1.14-gpu"
								       placeholder="quickdeploy/yao-tensorflow:1.14-gpu"/>
							</div>
							<div class="col-md-6">
								<label>CMD</label>
								<div class="form-group">
									<input type="text" class="form-control task-cmd" maxlength=""
									       placeholder="Command to bring up task"/>
								</div>
							</div>
							<div class="col-md-2">
								<label>Remove</label>
								<div class="form-group">
									<button type="button" class="btn btn-default task-remove">Remove</button>
								</div>
							</div>
							<div class="col-md-2">
								<label>Host Name</label>
								<div class="form-group">
									<input type="text" class="form-control task-name" maxlength="32"
									       placeholder="Task Name & Node Name" value="node1" required/>
								</div>
							</div>
							<div class="col-md-2">
								<label>Node Role<abbr title="Node role">?</abbr></label>
								<select class="form-control form-control task-is-ps" required>
									<option value="1">Parameter Server</option>
									<option value="0">Worker</option>
									<option value="0" selected>Default</option>
									<option value="0">Reducer</option>
								</select>
							</div>
							<div class="col-md-2">
								<label>CPU</label>
								<div class="form-group">
									<input type="number" class="form-control task-cpu" step="1" min="1" value="1"
									       placeholder="number of CPU required" required/>
								</div>
							</div>
							<div class="col-md-2">
								<label>Memory</label>
								<div class="form-group">
									<input type="number" class="form-control task-mem" step="1024" min="1024"
									       value="4096" placeholder="MB" required/>
								</div>
							</div>
							<div class="col-md-2 hidden">
								<label>GPU Model<abbr title="preferred GPU model">?</abbr></label>
								<select class="form-control form-control task-gpu-model" required>
									<option value="k40">K40</option>
									<option value="k80" selected>K80</option>
									<option value="P100">P100</option>
								</select>
							</div>
							<div class="col-md-2">
								<label>GPU</label>
								<div class="form-group">
									<input type="number" class="form-control task-gpu-num" step="1" min="1" value="1"
									       placeholder="number of GPU cards required" required/>
								</div>
							</div>
							<div class="col-md-2">
								<label>GPU Memory<abbr title="per card">?</abbr></label>
								<div class="form-group">
									<input type="number" class="form-control task-gpu-mem" step="1024" min="1024"
									       value="4096" placeholder="MB" required/>
								</div>
							</div>
						</div>
					</div>
					<div>
						<button id="form-job-submit" type="submit" class="btn btn-primary btn-lg">Submit</button>
						<button id="form-job-task-add" type="button" class="btn btn-default btn-lg">Add Task</button>
						<button id="form-job-predict-req" type="button" class="btn btn-default">Predict</button>
						<button id="form-job-predict-time" type="button" class="btn btn-default">PredictTime</button>
					</div>
				</form>
			</div>
		</div>
	</div>
</div>

<!-- virtual cluster modal -->
<div class="modal fade" id="modal-cluster" tabindex="-1" role="dialog" aria-labelledby="myModalLabel"
     aria-hidden="true">
	<div class="modal-dialog">
		<div class="modal-content panel-info">
			<div class="modal-header panel-heading">
				<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span>
				</button>
				<h4 id="modal-cluster-title" class="modal-title">Create Virtual Cluster</h4>
			</div>
			<div class="modal-body">
				<form class="form" action="javascript:void(0)">
					<label>Name</label>
					<div class="form-group form-group-lg">
						<label for="form-cluster-name" class="sr-only">Name</label>
						<input type="text" id="form-cluster-name" class="form-control" maxlength="16"
						       placeholder="virtual cluster name" required/>
					</div>
					<label>Reserved</label>
					<div class="form-group form-group-lg">
						<label for="form-cluster-reserved" class="sr-only">Reserved</label>
						<input type="checkbox" id="form-cluster-reserved"/>&nbsp;&nbsp;Reserved?
					</div>
					<label>Weight</label>
					<div class="form-group form-group-lg">
						<label for="form-cluster-weight" class="sr-only">Weight</label>
						<input type="number" id="form-cluster-weight" class="form-control" min="0" step="1"
						       value="10" placeholder=""/>
					</div>
					<label>GPU Number</label>
					<div class="form-group form-group-lg">
						<label for="form-cluster-quota-gpu-number" class="sr-only">GPU number</label>
						<input type="number" id="form-cluster-quota-gpu-number" class="form-control" min="0" step="1"
						       value="0" placeholder=""/>
					</div>
					<label>GPU Memory (Each)</label>
					<div class="form-group form-group-lg">
						<label for="form-cluster-quota-gpu-memory" class="sr-only">GPU Memory</label>
						<input type="number" id="form-cluster-quota-gpu-memory" class="form-control" min="0" step="1024"
						       value="10240" placeholder="(MB)"/>
					</div>
					<label>CPU</label>
					<div class="form-group form-group-lg">
						<label for="form-cluster-quota-cpu" class="sr-only">CPU</label>
						<input type="number" id="form-cluster-quota-cpu" class="form-control" placeholder="" min="0"
						       value="0" step="1"/>
					</div>
					<label>Memory (Total)</label>
					<div class="form-group form-group-lg">
						<label for="form-cluster-quota-mem" class="sr-only">CPU</label>
						<input type="number" id="form-cluster-quota-mem" class="form-control" placeholder="(MB)" min="0"
						       value="1024" step="1024"/>
					</div>
					<div>
						<input type="hidden" id="form-cluster-submit-type"/>
						<button id="form-cluster-submit" type="submit" class="btn btn-primary btn-lg">Submit</button>
					</div>
				</form>
			</div>
		</div>
	</div>
</div>

<!-- agent modal -->
<div class="modal fade" id="modal-agent" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
	<div class="modal-dialog">
		<div class="modal-content panel-info">
			<div class="modal-header panel-heading">
				<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span>
				</button>
				<h4 id="modal-agent-title" class="modal-title">Add New Agent</h4>
			</div>
			<div class="modal-body">
				<form class="form" action="javascript:void(0)">
					<label>IP</label>
					<div class="form-group form-group-lg">
						<label for="form-agent-ip" class="sr-only">IP</label>
						<input type="text" id="form-agent-ip" class="form-control" maxlength="64"
						       placeholder="10.0.0.1" required/>
					</div>
					<label>Alias</label>
					<div class="form-group form-group-lg">
						<label for="form-agent-alias" class="sr-only">Alias</label>
						<input type="text" id="form-agent-alias" class="form-control" maxlength="32"
						       placeholder="bj.node1"/>
					</div>
					<label>Cluster</label>
					<div class="form-group form-group-lg">
						<label for="form-agent-cluster" class="sr-only">Cluster</label>
						<select id="form-agent-cluster" class="form-control">
							<option value="0">default</option>
						</select>
					</div>
					<label>Token</label>
					<div class="form-group form-group-lg">
						<label for="form-agent-token" class="sr-only">Token</label>
						<input type="text" id="form-agent-token" class="form-control" placeholder="******" readonly/>
					</div>
					<div>
						<input type="hidden" id="form-agent-submit-type"/>
						<input type="hidden" id="form-agent-id"/>
						<button id="form-agent-submit" type="submit" class="btn btn-primary btn-lg">Submit</button>
					</div>
				</form>
			</div>
		</div>
	</div>
</div>

<!-- workspace modal -->
<div class="modal fade" id="modal-workspace" tabindex="-1" role="dialog" aria-labelledby="myModalLabel"
     aria-hidden="true">
	<div class="modal-dialog">
		<div class="modal-content panel-info">
			<div class="modal-header panel-heading">
				<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span>
				</button>
				<h4 id="modal-workspace-title" class="modal-title">Add New Workspace</h4>
			</div>
			<div class="modal-body">
				<form class="form" action="javascript:void(0)">
					<label>Name</label>
					<div class="form-group form-group-lg">
						<label for="form-agent-ip" class="sr-only">IP</label>
						<input type="text" id="form-workspace-name" class="form-control" maxlength="64"
						       placeholder="workspace name" required/>
					</div>
					<label>Type</label>
					<div class="form-group form-group-lg">
						<label for="form-workspace-type" class="sr-only">Type</label>
						<select id="form-workspace-type" class="form-control">
							<option value="git">git</option>
						</select>
					</div>
					<label>Git Repo</label>
					<div class="form-group form-group-lg">
						<label for="form-workspace-git-repo" class="sr-only">Git Repo</label>
						<input type="text" id="form-workspace-git-repo" class="form-control"
						       placeholder="http://192.168.100.100:3000/newnius/tf.git"/>
					</div>
					<div>
						<input type="hidden" id="form-workspace-submit-type"/>
						<input type="hidden" id="form-workspace-id"/>
						<input type="hidden" id="form-workspace-content"/>
						<button id="form-workspace-submit" type="submit" class="btn btn-primary btn-lg">Create</button>
					</div>
				</form>
			</div>
		</div>
	</div>
</div>
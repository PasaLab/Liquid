# YAO-scheduler


## API

#### ResourcePool
**GetHeartCounter**

```
?action=get_counter
```

**GetJobTaskStatusJHL**

```
?action=jhl_job_status&job=
```

#### Scheduler
**EnableSchedule**
```
?action=debug_enable
```

**DisableSchedule**
```
?action=debug_disable
```

**UpdateMaxParallelism**
```
?action=debug_update_parallelism&parallelism=5
```


**getAllPredicts**
```
?action=debug_get_predicts
```


**getAllGPUUtils**
```
?action=debug_get_gpu_utils
```


**SetShareRatio**
```
?action=debug_update_enable_share_ratio&ratio=0.75
```


**SetPreScheduleRatio**
```
?action=debug_update_enable_pre_schedule_ratio&ratio=0.95
```

**UpdateAllocateStrategy**
```
?action=allocator_update_strategy&strategy=bestfit
```

**SchedulerDump**
```
?action=debug_scheduler_dump
```

**DescribeJob**
```
?action=debug_optimizer_describe_job&job=
```

**EnableBatchAllocation**
```
?action=pool_enable_batch
```

**DisableBatchAllocation**
```
?action=pool_disable_batch
```

**UpdateBatchInterval**
```
?action=pool_set_batch_interval&interval=30
```

**PoolDump**
```
?action=debug_pool_dump
```

**EnableMock**
```
?action=debug_enable_mock
```

**DisableMock**
```
?action=debug_disable_mock
```

**UpdateStrategy**
```
?aciotn=allocator_update_strategy&strategy=mixed
```

**UpdateShareMaxUtilization**
```
?aciotn=conf_set_share_max_utilization&util=1.5
```
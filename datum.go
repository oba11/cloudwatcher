package main

import (
	"github.com/cloudfoundry/gosigar"
	"github.com/crowdmob/goamz/cloudwatch"
	"time"
)

func GetDimensions() []cloudwatch.Dimension {
	return []cloudwatch.Dimension{
		cloudwatch.Dimension{
			Name:  "InstanceId",
			Value: Config.InstanceId,
		},
	}
}

func GetFileSystemDatum() []cloudwatch.MetricDatum {
	metrics := []cloudwatch.MetricDatum{}

	dimensions := GetDimensions()

	now := time.Now()

	disk := sigar.FileSystemUsage{}
	disk.Get(Config.DiskPath)

	// Convert to bytes
	disk.Total = (disk.Total << 1) / 8 * 4096
	disk.Used = (disk.Used << 1) / 8 * 4096
	disk.Avail = (disk.Avail << 1) / 8 * 4096

	if Config.DiskSpaceUtil {
		metrics = append(metrics, cloudwatch.MetricDatum{
			Dimensions: dimensions,
			MetricName: "DiskSpaceUtilization",
			Unit:       "Percent",
			Timestamp:  now,
			Value:      float64(100 * disk.Used / disk.Total),
		})
	}

	if Config.DiskSpaceUsed {
		metrics = append(metrics, cloudwatch.MetricDatum{
			Dimensions: dimensions,
			MetricName: "DiskSpaceUsed",
			Unit:       Config.DiskSpaceUnits,
			Timestamp:  now,
			Value:      float64(disk.Used / Config.DiskSpaceUnitsInt()),
		})
	}

	if Config.DiskSpaceAvail {
		metrics = append(metrics, cloudwatch.MetricDatum{
			Dimensions: dimensions,
			MetricName: "DiskSpaceAvailable",
			Unit:       Config.DiskSpaceUnits,
			Timestamp:  now,
			Value:      float64(disk.Avail / Config.DiskSpaceUnitsInt()),
		})
	}

	return metrics
}

func GetMemoryDatum() []cloudwatch.MetricDatum {
	metrics := []cloudwatch.MetricDatum{}

	dimensions := GetDimensions()

	now := time.Now()

	mem := sigar.Mem{}
	mem.Get()

	swap := sigar.Swap{}
	swap.Get()

	if Config.MemUtil {
		metrics = append(metrics, cloudwatch.MetricDatum{
			Dimensions: dimensions,
			MetricName: "MemoryUtilization",
			Unit:       "Percent",
			Timestamp:  now,
			Value:      float64(100 * mem.Used / mem.Total),
		})
	}

	if Config.MemAvail {
		metrics = append(metrics, cloudwatch.MetricDatum{
			Dimensions: dimensions,
			MetricName: "MemoryAvailable",
			Unit:       Config.MemoryUnits,
			Timestamp:  now,
			Value:      float64(mem.Free / Config.MemoryUnitsInt()),
		})
	}

	if Config.MemUsed {
		metrics = append(metrics, cloudwatch.MetricDatum{
			Dimensions: dimensions,
			MetricName: "MemoryUsed",
			Unit:       Config.MemoryUnits,
			Timestamp:  now,
			Value:      float64(mem.Used / Config.MemoryUnitsInt()),
		})
	}

	if Config.SwapUtil {
		metrics = append(metrics, cloudwatch.MetricDatum{
			Dimensions: dimensions,
			MetricName: "SwapUtilization",
			Unit:       "Percent",
			Timestamp:  now,
			Value:      float64(100 * swap.Used / swap.Total),
		})
	}

	if Config.SwapUsed {
		metrics = append(metrics, cloudwatch.MetricDatum{
			Dimensions: dimensions,
			MetricName: "SwapUsed",
			Unit:       Config.MemoryUnits,
			Timestamp:  now,
			Value:      float64(swap.Used / Config.MemoryUnitsInt()),
		})
	}

	return metrics
}

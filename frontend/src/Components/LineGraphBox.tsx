import React, { useState, useEffect } from "react";
import { DeviceDropdown, MetricDropdown } from "./Dropdowns";
import { Device, MetricType } from "./Dropdowns";
import { GetFilter, MetricData } from "../../types/DataSnapshot";
import { filteredGetRequest } from "../App";
import LineChart from "./Linechart";

interface props {
  devices: Device[];
  metrictypes: MetricType[];
}

const LineGraphBox = ({ devices, metrictypes }: props) => {
  const [metricData, setMetricData] = useState<MetricData[] | undefined>(
    undefined
  );

  const [metricsFilters, setMetricsFilters] = useState<GetFilter>({
    device_id: 0,
    limit: 100,
    metric_id: 0,
    page_number: 1,
    setData: setMetricData, // This is for setting liveMetric state
  });

  // Fetch data based on the selected dropdown values
  useEffect(() => {
    filteredGetRequest(metricsFilters);
  }, [metricsFilters]);

  return (
    <div className="p-5">
      <h1 className="text-xl mb-4">Device Metrics Viewer</h1>

      <div className="flex gap-4 mb-5">
        {/* Device Dropdown */}
        <div className="w-1/3">
          <label className="block mb-2">Select a Device</label>
          <DeviceDropdown
            devices={devices}
            onDeviceSelect={(id) =>
              setMetricsFilters({ ...metricsFilters, device_id: id })
            }
          />
        </div>

        {/* Metric Dropdown */}
        <div className="w-1/3">
          <label className="block mb-2">Select a Metric</label>
          <MetricDropdown
            metrics={metrictypes}
            onMetricSelect={(id) =>
              setMetricsFilters({ ...metricsFilters, metric_id: id })
            }
          />
        </div>

        {/* Limit Input */}
        <div className="w-1/3">
          <label className="block mb-2">Set Data Limit</label>
          <input
            type="number"
            min="1"
            value={metricsFilters.limit || 100}
            onChange={(evt) => setMetricsFilters({ ...metricsFilters, limit: parseInt(evt.target.value, 10) || 100 })}
            className="w-full p-2 border border-gray-300 rounded"
            placeholder="Enter limit"
          />
        </div>
      </div>

      <div>
        <h2 className="text-lg mb-3">Metric Chart</h2>
        {metricData ? (
          <LineChart
            data={metricData.reverse()}
            metricName={
              metrictypes.find(
                (m: MetricType) => m.MetricID === metricsFilters.metric_id
              )?.Name || "Assortment of metrics"
            }
            chartTitle="Device Metric Over Time"
            yAxisLabel="Metric Value"
          />
        ) : (
          <p className="text-gray-500">
            Please select a device and metric to view the chart.
          </p>
        )}
      </div>
    </div>
  );
};

export default LineGraphBox;

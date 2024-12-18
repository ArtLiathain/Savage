import React, { useEffect, useState } from "react";
import { GaugeComponent } from "react-gauge-component";
import {
  GetFilter,
  MetricData,
  MetricType,
  Device,
} from "../../types/DataSnapshot";
import formatMetricName from "./utils";
import { filteredGetRequest, host, makeGetRequest } from "../App";
import { DeviceDropdown, MetricDropdown } from "./Dropdowns";

interface GaugeDisplayProps {
  devices: Device[];
  metricTypes: MetricType[];
}

const GaugeDisplay = ({ devices, metricTypes }: GaugeDisplayProps) => {
  const [liveMetric, setLiveMetrics] = useState<MetricData[] | undefined>(
    undefined
  );
  const [liveMetricsFilters, setLiveMetricsFilters] = useState<GetFilter>({
    device_id: 0,
    limit: 100,
    metric_id: 2,
    page_number: 1,
    setData: setLiveMetrics, // This is for setting liveMetric state
  });

  const calculateMaxLimit = (metricData: MetricData[] | undefined): number => {
    if (!metricData || metricData.length === 0) return 0;

    const metricName = metricData[0]?.MetricName || "";
    if (metricName.includes("percent")) return 100;

    return Math.max(...metricData.map((m) => m.Value));
  };

  useEffect(() => {
    filteredGetRequest(liveMetricsFilters);

    const interval = setInterval(() => {
      filteredGetRequest(liveMetricsFilters);
    }, 1000);

    return () => clearInterval(interval);
  }, [liveMetricsFilters]);
  let maxLimit: number | null = null;
  if (liveMetric) {
    maxLimit = calculateMaxLimit(liveMetric);
  }

  return (
    <div className="flex flex-col justify-center items-center w-full h-full max-w-full">
      <h1 className="text-xl mb-4">Gauge Display</h1>

      <div className="flex flex-wrap gap-4 mb-5 w-full max-w-5xl justify-between">
        {/* Device Dropdown */}
        <div className="flex-1 min-w-[200px]">
          <label className="block mb-2">Select a Device</label>
          <DeviceDropdown
            devices={devices}
            onDeviceSelect={(id) =>
              setLiveMetricsFilters({ ...liveMetricsFilters, device_id: id })
            }
          />
        </div>

        {/* Metric Dropdown */}
        <div className="flex-1 min-w-[200px]">
          <label className="block mb-2">Select a Metric</label>
          <MetricDropdown
            metrics={metricTypes}
            onMetricSelect={(id) =>
              setLiveMetricsFilters({ ...liveMetricsFilters, metric_id: id })
            }
          />
        </div>

        {/* Reset Button */}
        <div className="flex-1 min-w-[200px]">
          <label className="block mb-2">ESP low power</label>
          <button
            onClick={() => makeGetRequest(`${host}/reset`)}
            className="py-2 px-4 bg-red-500 text-white font-semibold rounded hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-400 focus:ring-opacity-50"
          >
            Reset
          </button>
        </div>
      </div>
      {liveMetric && (
        <>
          <p className="mb-3">{formatMetricName(liveMetric[0].MetricName)}</p>

          <GaugeComponent
            type="semicircle"
            arc={{
              colorArray: ["#00FF15", "#FF2121"],
              padding: 0.02,
              subArcs: [
                { limit: maxLimit! * 0.4 },
                { limit: maxLimit! * 0.6 },
                { limit: maxLimit! * 0.7 },
                { limit: maxLimit! },
              ],
            }}
            pointer={{ type: "blob", animationDelay: 0 }}
            value={liveMetric[0].Value}
            maxValue={calculateMaxLimit(liveMetric)}
            className="w-full max-w-md"
          />
        </>
      )}
    </div>
  );
};

export default GaugeDisplay;

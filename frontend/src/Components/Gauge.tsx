import React, { useEffect, useState } from "react";
import { GaugeComponent } from "react-gauge-component";
import { GetFilter, MetricData } from "../../types/DataSnapshot";
import formatMetricName from "./utils";
import { filteredGetRequest } from "../App";

// Define the GaugeComponent that accepts MetricData
const GaugeDisplay = () => {
  // Define the minimum and maximum values for the gauge (you can adjust these as needed)
  const [liveMetric, setLiveMetrics] = useState<MetricData[] | undefined>(
    undefined
  );
  const [liveMetricsFilters, setLiveMetricsFilters] = useState<GetFilter>({
    device_id: null,
    limit: 1,
    metric_id: 2,
    page_number: 1,
    setData: setLiveMetrics, // This is for setting liveMetric state
  });
  useEffect(() => {
    filteredGetRequest(liveMetricsFilters);

    const interval = setInterval(() => {
      filteredGetRequest(liveMetricsFilters);
    }, 1000);

    return () => clearInterval(interval);
  }, [liveMetricsFilters]);
  if (!liveMetric) {
    return <></>;
  }
  return (
    <div className="flex flex-col justify-center items-center w-fullh-full">
      <p>{formatMetricName(liveMetric[0].MetricName)}</p>

      <GaugeComponent
        type="semicircle"
        arc={{
          colorArray: ["#00FF15", "#FF2121"],
          padding: 0.02,
          subArcs: [
            { limit: 40 },
            { limit: 60 },
            { limit: 70 },
            {},
            {},
            {},
            {},
          ],
        }}
        pointer={{ type: "blob", animationDelay: 0 }}
        value={liveMetric[0].Value}
      />
    </div>
  );
};
export default GaugeDisplay;

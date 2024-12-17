import { GetFilter, MetricData } from "../types/DataSnapshot";
import Card from "./Components/Card";
import { Dispatch, SetStateAction, useEffect, useState } from "react";
import LineChart from "./Components/Linechart";
import MetricDataTable from "./Components/MetricTable";
import GaugeDisplay from "./Components/Gauge";

function App() {
  const [liveMetrics, setLiveMetrics] = useState<MetricData[] | undefined>(
    undefined
  );
  const [liveMetricsfilters, setLiveMetricsFilters] = useState<GetFilter>({
    device_id: null,
    limit: 1,
    metric_id: 2,
    page_number: 1,
    setData: setLiveMetrics,
  });
  const [metrics, setMetrics] = useState<MetricData[] | undefined>(undefined);

  const filteredGetRequest = async ({
    metric_id,
    limit,
    page_number,
    device_id,
    setData,
  }: GetFilter) => {
    let url = `http://localhost:8080/read?metric_id=${metric_id}&limit=${limit}&page_number=${page_number}&device_id=${device_id}`;
    try {
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error("Failed to fetch data");
      }
      const data = await response.json();
      setData(data.reverse());
    } catch (error) {
      console.error("Error fetching metrics:", error);
    }
  };
  useEffect(() => {
    filteredGetRequest(liveMetricsfilters);

    const interval = setInterval(filteredGetRequest, 1000);

    return () => clearInterval(interval);
  }, []);

  return (
    <div className="bg-green-500 h-screen grid grid-cols-2 gap-4 p-10">
      {/* Ensure grid items respect container height */}
      <Card>
        {metrics && (
          <LineChart
            data={metrics}
            metricName="ram_usage_percent"
            chartTitle="RAM Usage Over Time"
            yAxisLabel="RAM Usage (%)"
            numberOfDataPoints={50}
          />
        )}
      </Card>
      <Card>{metrics && <MetricDataTable data={metrics} />}</Card>
      <Card>{liveMetrics && <GaugeDisplay metric={liveMetrics[0]} />}</Card>
      <Card>{metrics && <MetricDataTable data={metrics} />}</Card>
    </div>
  );
}

export default App;

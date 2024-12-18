import {
  Device,
  GetFilter,
  MetricData,
  MetricType,
} from "../types/DataSnapshot";
import Card from "./Components/Card";
import { useEffect, useState } from "react";
import LineChart from "./Components/Linechart";
import MetricDataTable from "./Components/MetricTable";
import GaugeDisplay from "./Components/Gauge";

const host = "http://localhost:8080";

export const filteredGetRequest = async ({
  metric_id,
  limit,
  page_number,
  device_id,
  setData,
}: GetFilter) => {
  let url = `${host}/read?metric_id=${metric_id}&limit=${limit}&page_number=${page_number}&device_id=${device_id}`;
  try {
    let snapshots = await makeGetRequest(url);
    setData(snapshots.reverse());
  } catch {
    console.log("Failed getting snapshots");
  }
};

export const makeGetRequest = async (url: string) => {
  try {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error("Failed to fetch data");
    }
    const data = await response.json();
    return data; // Store the reversed data
  } catch (error) {
    throw new Error("Failed to fetch data");
  }
};

function App() {
  const [metrics, setMetrics] = useState<MetricData[] | undefined>(undefined);
  const [devices, setDevices] = useState<Device[] | undefined>(undefined);
  const [metricTypes, setMetricTypes] = useState<MetricType[] | undefined>(
    undefined
  );

  // This function performs a GET request with filters to fetch metrics data

  useEffect(() => {
    const setconstdata = async () => {
      try {
        let devices = await makeGetRequest(`${host}/devices`);
        setDevices(devices);
      } catch {
        console.log("Failed getting snapshots");
      }
      try {
        let metricTypes = await makeGetRequest(`${host}/metrictypes`);
        setMetricTypes(metricTypes);
      } catch {
        console.log("Failed getting snapshots");
      }
    };
    setconstdata();
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
      <Card>
        <GaugeDisplay />
      </Card>
      <Card>{metrics && <MetricDataTable data={metrics} />}</Card>
    </div>
  );
}

export default App;

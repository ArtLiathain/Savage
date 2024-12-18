import {
  Device,
  GetFilter,
  MetricData,
  MetricType,
} from "../types/DataSnapshot";
import Card from "./Components/Card";
import { useEffect, useState } from "react";
import MetricDataTable from "./Components/MetricTable";
import GaugeDisplay from "./Components/Gauge";
import LineGraphBox from "./Components/LineGraphBox";

export const host = "http://209.97.179.7:8080";

export const filteredGetRequest = async ({
  metric_id,
  limit,
  page_number,
  device_id,
  setData,
}: GetFilter) => {
  let url = `${host}/read?metric_id=${metric_id}&limit=${limit}&page=${page_number}&device_id=${device_id}`;
  console.log(url)
  try {
    let snapshots = await makeGetRequest(url);
    setData(snapshots);
  } catch (err) {
    console.log(err);
    console.log("Failed getting snapshots");
  }
};

export const makeGetRequest = async (url: string) => {
  try {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error("Failed to fetch at all");
    }
    const data = await response.json();
    return data; // Store the reversed data
  } catch (error) {
    throw new Error("Failed to fetch data");
  }
};

function App() {
  const [devices, setDevices] = useState<Device[]>([]);
  const [metricTypes, setMetricTypes] = useState<MetricType[]>([]);

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
    <div className="bg-green-500 min-h-screen p-10 overflow-auto grid grid-rows-[auto,1fr,auto] gap-4">
      {/* Line Graph Section */}
      <div className="row-span-1 col-span-1 w-full">
        <Card>
          <LineGraphBox devices={devices} metrictypes={metricTypes} />
        </Card>
      </div>

      {/* Table Section */}
      <div className="row-span-1 col-span-1 w-full">
        <Card>
          <MetricDataTable devices={devices} metricTypes={metricTypes} />
        </Card>
      </div>

      {/* Gauge Section */}
      <div className="row-span-1 col-span-1 w-full">
        <Card>
          <GaugeDisplay devices={devices} metricTypes={metricTypes} />
        </Card>
      </div>
    </div>
  );
}
export default App;

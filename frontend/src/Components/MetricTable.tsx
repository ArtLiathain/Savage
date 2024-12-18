import { useState, useEffect } from "react";
import {
  MetricData,
  Device,
  MetricType,
  GetFilter,
} from "../../types/DataSnapshot";
import { filteredGetRequest } from "../App";
import { DeviceDropdown, MetricDropdown } from "./Dropdowns";

interface DataTableProps {
  devices: Device[];
  metricTypes: MetricType[];
}

const MetricDataTable = ({ devices, metricTypes }: DataTableProps) => {
  const [metricData, setMetricData] = useState<MetricData[] | undefined>(
    undefined
  );
  const [filters, setFilters] = useState<GetFilter>({
    device_id: 0,
    metric_id: 0,
    limit: 10,
    page_number: 1,
    setData: setMetricData,
  });
  
  // Fetch data based on selected filters
  useEffect(() => {
    filteredGetRequest(filters);
  }, [filters]);

  // Handle device selection
  const handleDeviceSelect = (deviceId: number | null) => {
    setFilters((prevFilters) => ({
      ...prevFilters,
      device_id: deviceId,
      page_number : 1
    }));
  };

  const handleMetricSelect = (metricId: number | null) => {
    setFilters((prevFilters) => ({
      ...prevFilters,
      metric_id: metricId,
      page_number : 1
    }));
  };

  // Pagination functionality
  const handleNextPage = () => {
    setFilters((prevFilters) => ({
      ...prevFilters,
      page_number: (prevFilters.page_number ?? 1) + 1,
    }));
  };

  const handlePreviousPage = () => {
    setFilters((prevFilters) => ({
      ...prevFilters,
      page_number: Math.max((prevFilters.page_number ?? 1) - 1, 1), // Ensure the page doesn't go below 1
    }));
  };

  return (
    <div className="flex flex-col items-center w-full p-5">
      <h1 className="text-xl mb-4">Metric Data Viewer</h1>

      {/* Dropdowns for device and metric */}
      <div className="flex gap-20 mb-5">
        <div className="w-full sm:w-1/3">
          <label className="block mb-2">Select a Device</label>
          <DeviceDropdown
            devices={devices}
            onDeviceSelect={handleDeviceSelect}
          />
        </div>

        <div className="w-full sm:w-1/3">
          <label className="block mb-2">Select a Metric</label>
          <MetricDropdown
            metrics={metricTypes}
            onMetricSelect={handleMetricSelect}
          />
        </div>
      </div>

      {/* Data Table */}
      {metricData && metricData.length > 0 ? (
        <table className="min-w-full table-auto border-collapse">
          <thead className="bg-gray-100">
            <tr>
              <th className="px-4 py-2 text-left">Metric Name</th>
              <th className="px-4 py-2 text-left">Device Name</th>
              <th className="px-4 py-2 text-left">Value</th>
              <th className="px-4 py-2 text-left">Client Time</th>
              <th className="px-4 py-2 text-left">Server Time</th>
            </tr>
          </thead>
          <tbody>
            {metricData.map((data, index) => (
              <tr key={index} className="border-b hover:bg-gray-50">
                <td className="px-4 py-2">{data.MetricName}</td>
                <td className="px-4 py-2">{data.DeviceName}</td>
                <td className="px-4 py-2">{data.Value}</td>
                <td className="px-4 py-2">{data.ClientUtcTime}</td>
                <td className="px-4 py-2">{data.ServerUtcTime}</td>
              </tr>
            ))}
          </tbody>
        </table>
      ) : (
        <p className="text-gray-500 mt-4">
          No data available. Please select a device and metric.
        </p>
      )}

      {/* Pagination Controls */}
      <div className="flex justify-center items-center mt-4">
        <button
          onClick={handlePreviousPage}
          className="px-4 py-2 bg-blue-500 text-white rounded-l hover:bg-blue-600"
        >
          Previous
        </button>
        <span className="px-4 py-2">Page {filters.page_number}</span>
        <button
          onClick={handleNextPage}
          className="px-4 py-2 bg-blue-500 text-white rounded-r hover:bg-blue-600"
        >
          Next
        </button>
      </div>
    </div>
  );
};
export default MetricDataTable;

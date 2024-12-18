import React from "react";

export interface MetricType {
  MetricID: number;
  Name: string;
}

export interface Device {
  DeviceName: string;
  DeviceGuid: string;
  DeviceID: number;
}

interface DropdownProps<T> {
  data: T[]; // Generic array of data
  labelKey: keyof T; // Key to display in the dropdown (like "Name")
  idKey: keyof T; // Key to extract the ID (like "DeviceID")
  onSelect: (id: number) => void; // Callback function
  placeholder: string; // Placeholder text for the dropdown
}

function Dropdown<T extends object>({
  data,
  labelKey,
  idKey,
  onSelect,
  placeholder,
}: DropdownProps<T>) {
  const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const selectedId = parseInt(e.target.value, 10); // Parse the ID
    onSelect(selectedId); // Call the callback
  };

  return (
    <select onChange={handleChange} className="p-2 border rounded">
      <option value={0}>{placeholder}</option>
      {data.map((item, index) => (
        <option key={index} value={item[idKey] as unknown as number}>
          {item[labelKey] as string}
        </option>
      ))}
    </select>
  );
}

// Export specialized components for Devices and Metrics
export const DeviceDropdown: React.FC<{
  devices: Device[];
  onDeviceSelect: (id: number) => void;
}> = ({ devices, onDeviceSelect }) => (
  <Dropdown
    data={devices}
    labelKey="DeviceName"
    idKey="DeviceID"
    onSelect={onDeviceSelect}
    placeholder="Select a Device"
  />
);

export const MetricDropdown: React.FC<{
  metrics: MetricType[];
  onMetricSelect: (id: number) => void;
}> = ({ metrics, onMetricSelect }) => (
  <Dropdown
    data={metrics}
    labelKey="Name"
    idKey="MetricID"
    onSelect={onMetricSelect}
    placeholder="Select a Metric"
  />
);

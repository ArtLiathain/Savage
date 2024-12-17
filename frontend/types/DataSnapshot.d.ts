export interface MetricData {
  Value: number;
  MetricName: string;
  DeviceGuid: string;
  DeviceName: string;
  ClientUtcTime: string;
  ClientTimezoneMinutes: number;
  ServerUtcTime: string;
  ServerTimezoneMinutes: number;
}

export interface GetFilter {
  device_id: number | null;
  metric_id: number | null;
  limit: number | null;
  page_number: number | null;
  device_id: number | null;
  setData: Dispatch<SetStateAction<MetricData | undefined>>;
}

export interface Metric {
    Name: string;
    Value: number;
}

export interface DataSnapshot {
    Id: string;
    Timestamp: string;
    Metrics: Metric[];
}
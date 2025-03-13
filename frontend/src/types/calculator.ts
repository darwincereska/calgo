export interface CalculatorValue {
    number: number;
    unit: string;
    type: ValueType;
    isError: boolean;
    errorMsg?: string;
}

export enum ValueType {
    Number,
    Currency,
    Unit,
    Time,
    Percentage,
    Date,
    Temperature
}

export interface CalculatorContext {
    variables: Map<string, CalculatorValue>;
    timezone: string;
    ppi: number;
    emSize: number;
}
export interface GrowUnit {
  Id: number;
  Name: string;
  Width: number;
  Height: number;
  Depth: number;
  CarbonFilter: boolean;
  ActiveIntake: boolean;
  OuttakeFanThroughputInM3H: number;
  WattageLamp: number;
  Ventilation: boolean;
  Inside: boolean;
}

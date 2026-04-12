export interface TagData {
  uid: string;
  date: string;
  dateDisplay: string;
  supplierCode: string;
  supplierName: string;
  materialCode: string;
  materialName: string;
  color: string;
  lengthCode: string;
  lengthDisplay: string;
  serial: string;
}

export interface WriteRequest {
  date: string;
  supplier: string;
  material: string;
  color: string;
  length: string;
  serial: string;
}

export interface MaterialOption {
  code: string;
  name: string;
}

export interface VendorOption {
  code: string;
  name: string;
}

export interface LengthOption {
  code: string;
  name: string;
  grams: string;
}

export interface OptionsResponse {
  materials: MaterialOption[];
  vendors: VendorOption[];
  lengths: LengthOption[];
}

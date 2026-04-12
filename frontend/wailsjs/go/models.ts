export namespace main {
	
	export class LengthOption {
	    code: string;
	    name: string;
	    grams: string;
	
	    static createFrom(source: any = {}) {
	        return new LengthOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.name = source["name"];
	        this.grams = source["grams"];
	    }
	}
	export class MaterialOption {
	    code: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new MaterialOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.name = source["name"];
	    }
	}
	export class VendorOption {
	    code: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new VendorOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.name = source["name"];
	    }
	}
	export class OptionsResponse {
	    materials: MaterialOption[];
	    vendors: VendorOption[];
	    lengths: LengthOption[];
	
	    static createFrom(source: any = {}) {
	        return new OptionsResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.materials = this.convertValues(source["materials"], MaterialOption);
	        this.vendors = this.convertValues(source["vendors"], VendorOption);
	        this.lengths = this.convertValues(source["lengths"], LengthOption);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TagData {
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
	
	    static createFrom(source: any = {}) {
	        return new TagData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uid = source["uid"];
	        this.date = source["date"];
	        this.dateDisplay = source["dateDisplay"];
	        this.supplierCode = source["supplierCode"];
	        this.supplierName = source["supplierName"];
	        this.materialCode = source["materialCode"];
	        this.materialName = source["materialName"];
	        this.color = source["color"];
	        this.lengthCode = source["lengthCode"];
	        this.lengthDisplay = source["lengthDisplay"];
	        this.serial = source["serial"];
	    }
	}
	
	export class WriteRequest {
	    date: string;
	    supplier: string;
	    material: string;
	    color: string;
	    length: string;
	    serial: string;
	
	    static createFrom(source: any = {}) {
	        return new WriteRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.supplier = source["supplier"];
	        this.material = source["material"];
	        this.color = source["color"];
	        this.length = source["length"];
	        this.serial = source["serial"];
	    }
	}

}


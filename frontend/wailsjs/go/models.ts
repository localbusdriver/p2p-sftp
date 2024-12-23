export namespace handlers {
	
	export class FileUpload {
	    Id: string;
	    Filename: string;
	    StoragePath: string;
	    Size: number;
	    // Go type: time
	    Uploadtime: any;
	    UserID: string;
	    Status: string;
	
	    static createFrom(source: any = {}) {
	        return new FileUpload(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Filename = source["Filename"];
	        this.StoragePath = source["StoragePath"];
	        this.Size = source["Size"];
	        this.Uploadtime = this.convertValues(source["Uploadtime"], null);
	        this.UserID = source["UserID"];
	        this.Status = source["Status"];
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

}


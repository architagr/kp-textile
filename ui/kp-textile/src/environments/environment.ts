// This file can be replaced during build by using the `fileReplacements` array.
// `ng build` replaces `environment.ts` with `environment.prod.ts`.
// The list of file replacements can be found in `angular.json`.

export const environment = {
  production: false,
  clientBaseUrl: "http://localhost:8080/",
  qualityBaseUrl: "http://localhost:8080/",
  hsnBaseUrl: "http://localhost:8081/hsncode/",
  vendorBaseUrl: "http://localhost:8082/",
  transporterBaseUrl: "http://localhost:8083/transporter/",
  salesBaseUrl: "http://localhost:8084/sales",
  bailBaseUrl: "http://localhost:8084/bailInfo",
  purchaseBaseUrl: "http://localhost:8084/purchase",
  documentBaseUrl: "http://localhost:8085/",
  organizationBaseUrl: "http://localhost:8087/",
};

/*
 * For easier debugging in development mode, you can import the following file
 * to ignore zone related error stack frames such as `zone.run`, `zoneDelegate.invokeTask`.
 *
 * This import should be commented out in production mode because it will have a negative impact
 * on performance if an error is thrown.
 */
// import 'zone.js/plugins/zone-error';  // Included with Angular CLI.

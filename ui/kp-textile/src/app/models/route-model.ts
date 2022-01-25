export declare interface RouteInfo {
  path: string;
  title: string;
  icon: string;
  class: string;
  parentPath: string;
  childRoutes:string[];
  showMenu: boolean
}
export const ROUTES: RouteInfo[] = [
  {
    path: "/dashboard",
    title: "Dashboard",
    icon: "icon-chart-pie-36",
    class: "",
    parentPath:"",
    childRoutes:[],
    showMenu:true
  },
  {
    path: "/client",
    title: "Client List",
    icon: "icon-bank",
    class: "",
    parentPath:"",
    childRoutes:['updateclient','addclient'],
    showMenu:true
  },
  {
    path: "/addclient",
    title: "Add new Client",
    icon: "icon-bank",
    class: "",
    parentPath:"client",
    childRoutes:[],
    showMenu:false
  },
  {
    path: "/updateclient",
    title: "Update Client",
    icon: "icon-bank",
    class: "",
    parentPath:"client",
    childRoutes:[],
    showMenu:false
  },
  {
    path: "/vendor",
    title: "Vendor",
    icon: "icon-bus-front-12",
    class: "",
    parentPath:"",
    childRoutes:['addvendor','updatevendor'],
    showMenu:true
  },
  {
    path: "/addvendor",
    title: "Add new vendor",
    icon: "icon-bus-front-12",
    class: "",
    parentPath:"vendor",
    childRoutes:[],
    showMenu:false
  },
  {
    path: "/updatevendor",
    title: "Update Vendor",
    icon: "icon-bus-front-12",
    class: "",
    parentPath:"vendor",
    childRoutes:[],
    showMenu:false
  },
  {
    path: "/transpoter",
    title: "Transporters",
    icon: "icon-delivery-fast",
    class: "",
    parentPath:"",
    childRoutes:['addtranspoter','updatetranspoter'],
    showMenu:true
  },
  {
    path: "/addtranspoter",
    title: "Add new Transporter",
    icon: "icon-delivery-fast",
    class: "",
    parentPath:"transpoter",
    childRoutes:[],
    showMenu:false
  },
  {
    path: "/updatetranspoter",
    title: "Update Transporter",
    icon: "icon-delivery-fast",
    class: "",
    parentPath:"transpoter",
    childRoutes:[],
    showMenu:false
  },
  {
    path: "/quality",
    title: "Quality",
    icon: "icon-scissors",
    class: "",
    parentPath:"",
    childRoutes:[],
    showMenu:true
  },
  {
    path: "/hsncode",
    title: "HSN Code",
    icon: "icon-molecule-40",
    class: "",
    parentPath:"",
    childRoutes:[],
    showMenu:true
  },
  {
    path: "/purchase",
    title: "Purcahse",
    icon: "icon-cart",
    class: "",
    parentPath:"",
    childRoutes:['addpurchase','updatepurchase'],
    showMenu:true
  },
  {
    path: "/addpurchase",
    title: "Add Purcahse",
    icon: "icon-cart",
    class: "",
    parentPath:"purchase",
    childRoutes:[],
    showMenu:false
  },
  {
    path: "/updatepurchase",
    title: "Update Purcahse",
    icon: "icon-cart",
    class: "",
    parentPath:"purchase",
    childRoutes:[],
    showMenu:false
  },
  {
    path: "/sales",
    title: "Sales",
    icon: "icon-coins",
    class: "",
    parentPath:"",
    childRoutes:['updatesales', 'addsales'],
    showMenu:true
  },
  {
    path: "/updatesales",
    title: "Update Sales Order",
    icon: "icon-coins",
    class: "",
    parentPath:"sales",
    childRoutes:[],
    showMenu:false
  },
  {
    path: "/addsales",
    title: "Add Sales Order",
    icon: "icon-coins",
    class: "",
    parentPath:"sales",
    childRoutes:[],
    showMenu:false
  }
];

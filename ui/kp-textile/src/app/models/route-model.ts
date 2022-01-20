declare interface RouteInfo {
  path: string;
  title: string;
  icon: string;
  class: string;
  showMenu: boolean
}
export const ROUTES: RouteInfo[] = [
  {
    path: "/dashboard",
    title: "Dashboard",
    icon: "icon-chart-pie-36",
    class: "",
    showMenu:true
  },
  {
    path: "/client",
    title: "Client List",
    icon: "icon-bank",
    class: "",
    showMenu:true
  },
  {
    path: "/addclient",
    title: "Add new Client",
    icon: "icon-bank",
    class: "",
    showMenu:false
  },
  {
    path: "/updateclient",
    title: "Update Client",
    icon: "icon-bank",
    class: "",
    showMenu:false
  },
  {
    path: "/vendor",
    title: "Vendor",
    icon: "icon-bus-front-12",
    class: "",
    showMenu:true
  },
  {
    path: "/addvendor",
    title: "Add new vendor",
    icon: "icon-bank",
    class: "",
    showMenu:false
  },
  {
    path: "/updatevendor",
    title: "Update Vendor",
    icon: "icon-bank",
    class: "",
    showMenu:false
  },
  {
    path: "/transpoter",
    title: "Transporters",
    icon: "icon-delivery-fast",
    class: "",
    showMenu:true
  },
  {
    path: "/quality",
    title: "Quality",
    icon: "icon-scissors",
    class: "",
    showMenu:true
  },
  {
    path: "/hsncode",
    title: "HSN Code",
    icon: "icon-molecule-40",
    class: "",
    showMenu:true
  },
  {
    path: "/purchase",
    title: "Purcahse",
    icon: "icon-cart",
    class: "",
    showMenu:true
  },

  {
    path: "/sales",
    title: "Sales",
    icon: "icon-coins",
    class: "",
    showMenu:true
  }
];

const escpos = require("escpos");
escpos.Network = require("escpos-network");
const moment = require("moment");

// Replace with your printer's IP and Port
const device = new escpos.Network("192.168.1.148", 9100);

// Set the printer adapter
const printer = new escpos.Printer(device);

const data = require("./test.json");
const orders = data.order;
let item = 0;

const maxlineLength = 48;

const datetime = moment(data.bill_date);
console.log(datetime.format("DD/MM/YYYY"));

const formatter = new Intl.NumberFormat("th-TH", {
  style: "currency",
  currency: "THB", // Change to the desired currency code
});

// orders.forEach((order) => {
//   // printer.align("lt");
//   // printer.text(`${order.quantity} ${order.name} (${order.size})`);
//   // printer.align("rt");
//   // printer.text(`${order.price.toFixed(2)}\n`);

//   let orderLine = `${order.quantity} ${order.name} (${order.size})`;
//   let price = `${formatter.format(order.price.toFixed(2))}`;
//   let totalLine = orderLine.length + price.length;
//   let spacesNeeded = Math.max(0, maxlineLength - totalLine);
//   let spaces = " ".repeat(spacesNeeded);
//   let orderLineData = orderLine + spaces + price;
//   orderLineData = orderLineData.replace(/[฿]/g, "");
//   console.log(orderLineData);

//   item += order.quantity;
// });

device.open(() => {
  printer
    .font("b")
    .align("ct")
    .style("b")
    .size(1, 1)
    .text("MAX WALLET\n")
    .size(0.5, 0.5)
    .font("a")
    .text("Receipt\n")
    .align("ct")
    .text("---------------------------------------\n")
    // head bill
    .align("lt")
    .text(`TABLE: ${data.table}`);
  printer
    .align("lt")
    .text(`CASHIER: ${data.operator}`)
    .text(`CUSTOMER: ${data.customer.aka}`);

  // datetime
  let dateLine = `DATE: ${datetime.format("DD/MM/YYYY")}`;
  let timeLine = `TIME: ${datetime.format("HH:mm")}`;
  let totalDatetimeLine = dateLine.length + timeLine.length;
  let spacesNeededDateTime = Math.max(0, maxlineLength - totalDatetimeLine);
  let spacesDatetime = " ".repeat(spacesNeededDateTime);
  printer.align("lt").text(dateLine + spacesDatetime + timeLine + "");

  printer.align("ct").text("---------------------------------------\n");

  orders.forEach((order) => {
    let orderLine = `${order.quantity}  ${order.name} (${order.size})`;
    let price = `${formatter.format(order.price.toFixed(2))}`;
    let totalLine = orderLine.length + price.length;
    // Calculate the spaces needed to align the price to the right
    let spacesNeeded = Math.max(0, maxlineLength - totalLine); // Ensure spacesNeeded is not negative
    let spaces = " ".repeat(spacesNeeded);
    let orderLineData = orderLine + spaces + price;
    orderLineData = orderLineData.replace(/[฿]/g, "");
    printer.align("lt");
    printer.text(orderLineData + "\n");

    item += order.quantity;
  });

  printer
    .align("ct")
    .text("---------------------------------------\n")
    // footer bill
    .align("lt")
    .text(`ITEMS: ${item}\n`);

  // summary bill
  let textLine = `Subtotal: `;
  let dataLine = `${formatter.format(data.price.toFixed(2))}`;
  let totalLine = textLine.length + dataLine.length;
  let spacesNeeded = Math.max(0, 24 - totalLine);
  let spaces = " ".repeat(spacesNeeded);
  let textData = textLine + spaces + dataLine;
  textData = textData.replace(/[฿]/g, "");
  printer.align("rt").text(textData + "\n");

  textLine = `Discount: `;
  dataLine = `${formatter.format(data.discount.toFixed(2))}`;
  totalLine = textLine.length + dataLine.length;
  spacesNeeded = Math.max(0, 24 - totalLine);
  spaces = " ".repeat(spacesNeeded);
  textData = textLine + spaces + dataLine;
  textData = textData.replace(/[฿]/g, "");
  printer.align("rt").text(textData + "\n");

  textLine = `Service Charge(${data.percent_service_charge}%): `;
  dataLine = `${formatter.format(data.service_charge.toFixed(2))}`;
  totalLine = textLine.length + dataLine.length;
  spacesNeeded = Math.max(0, 24 - totalLine);
  spaces = " ".repeat(spacesNeeded);
  textData = textLine + spaces + dataLine;
  textData = textData.replace(/[฿]/g, "");
  printer.align("rt").text(textData + "\n");

  let grand_total = data.price_with_discount + data.service_charge;
  textLine = `Before VAT: `;
  dataLine = `${formatter.format(grand_total.toFixed(2))}`;
  totalLine = textLine.length + dataLine.length;
  spacesNeeded = Math.max(0, 24 - totalLine);
  spaces = " ".repeat(spacesNeeded);
  textData = textLine + spaces + dataLine;
  textData = textData.replace(/[฿]/g, "");
  printer.align("rt").text(textData + "\n");

  textLine = `VAT(7%): `;
  dataLine = `${formatter.format(data.vat.toFixed(2))}`;
  totalLine = textLine.length + dataLine.length;
  spacesNeeded = Math.max(0, 24 - totalLine);
  spaces = " ".repeat(spacesNeeded);
  textData = textLine + spaces + dataLine;
  textData = textData.replace(/[฿]/g, "");
  printer.align("rt").text(textData + "\n");

  textLine = `Rounding: `;
  dataLine = `${formatter.format(data.rounding.toFixed(2))}`;
  totalLine = textLine.length + dataLine.length;
  spacesNeeded = Math.max(0, 24 - totalLine);
  spaces = " ".repeat(spacesNeeded);
  textData = textLine + spaces + dataLine;
  textData = textData.replace(/[฿]/g, "");
  printer.align("rt").text(textData + "\n");

  printer.align("rt").text("=======================");
  textLine = `Total: `;
  dataLine = `${formatter.format(data.total.toFixed(2))}`;
  totalLine = textLine.length + dataLine.length;
  spacesNeeded = Math.max(0, 24 - totalLine);
  spaces = " ".repeat(spacesNeeded);
  textData = textLine + spaces + dataLine;
  textData = textData.replace(/[฿]/g, "");
  // printer.style("b");
  printer.align("rt").text(textData + "");
  printer.align("rt").text("=======================");

  printer.align("ct").text("---------------------------------------\n");
  // printer.style("b");
  printer.align("ct").text("Thank you for your order!\n").cut().close();
});

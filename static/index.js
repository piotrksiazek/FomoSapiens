const priceChart = document.getElementById("priceChart").getContext("2d");
const sentimentChart = document
  .getElementById("sentimentChart")
  .getContext("2d");

const defaultNumverOfDaysToFetch = 180;

const usdPrice = [];
const sentiment = [];
const days = [];

// for (let i = 0; i < defaultNumverOfDaysToFetch; i++) {
//   days.push(i.toString());
// }
const formatter = new Intl.DateTimeFormat("pl");
axios
  .get(`getDays/${defaultNumverOfDaysToFetch}`)
  .then((response) => {
    response.data.forEach((day) => {
      usdPrice.push(day.realprice);
      sentiment.push(day.sentiment);

      date = Date.parse(day.creationday);
      days.push(formatter.format(date));
    });
    var chart1 = new Chart(priceChart, {
      type: "line",
      data: {
        labels: days,
        datasets: [
          {
            label: "BTC price in USD",
            yAxisID: "price",
            data: usdPrice,
            fill: false,
            borderColor: "gold",
            tension: 0.1,
          },
          {
            label: "Reddit sentiment",
            yAxisID: "sentiment",
            data: sentiment,
            fill: false,
            borderColor: "red",
            tension: 0.1,
          },
        ],
      },
      options: {
        responsive: true, // Instruct chart js to respond nicely.
        maintainAspectRatio: false, // Add to prevent default behaviour of full-width/height
        scales: {
          sentiment: {
            type: "linear",
            display: true,
            position: "left",
            max: 50,
            min: 30,
          },
          price: {
            type: "linear",
            display: true,
            position: "right",

            // grid line settings
            grid: {
              drawOnChartArea: false, // only want the grid lines for one axis to show up
            },
          },
        },
      },
    });
  })
  .catch((error) => console.log(error));

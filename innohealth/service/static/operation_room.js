$(function () {
    $.getJSON('http://localhost:1323/operation/byroom?dbname=testdb', function (data) {
        console.log(data)
        Highcharts.chart('container', {

            chart: {
                type: 'heatmap',
                marginTop: 40,
                marginBottom: 80,
                plotBorderWidth: 1
            },


            title: {
                text: '수술방 수술 현황'
            },

            xAxis: {
                categories: data['xCategory']
            },

            yAxis: {
                categories: data['yCategory'],
                title: '시간'
            },

            colorAxis: {
                min: 0,
                minColor: '#FFFFFF',
                maxColor: Highcharts.getOptions().colors[1]
            },

            legend: {
                align: 'right',
                layout: 'vertical',
                margin: 0,
                verticalAlign: 'top',
                y: 25,
                symbolHeight: 280
            },

            tooltip: {
                formatter: function () {
                    return '<b>' + this.series.xAxis.categories[this.point.x] + '</b> <br><b>'
                        + this.series.yAxis.categories[this.point.y] + '</b><br><b>' +this.point.value + '</b>';
                }
            },

            series: [{
                name: '시간별 운용',
                borderWidth: 0,
                data: data['data'],
                dataLabels: {
                    enabled: true,
                }
            }],
            credits: {
                enabled: false
            },

            events: {
                click: function (e) {
                    console.log(e);
                }
            }

        });
    });
});

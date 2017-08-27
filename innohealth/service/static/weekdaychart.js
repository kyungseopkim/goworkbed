$(function () {    
    $.getJSON('http://localhost:1323/operation/byweekday?dbname=testdb', function (data) {
        //weekday = ['일', '월', '화', '수', '목', '금', '토']; 
        weekday = ['일', '토', '금', '목', '수', '화', '월']; 
        department = [];
        deptmap = {};
        chartdata = [];
        $.each(data, function (index, val) {
            dep = val['department'];

            if (dep in deptmap) {
                index = deptmap[dep];
            } else {
                index = department.length
                department.push(dep)
                deptmap[dep] = index
            }

            dow = val['weekday'];
            chartdata.push([index, weekday.indexOf(dow), val['count']]);
        });

        Highcharts.chart('container', {

            chart: {
                type: 'heatmap',
                marginTop: 40,
                marginBottom: 80,
                plotBorderWidth: 1
            },


            title: {
                text: '요일별 수술 현황'
            },

            xAxis: {
                categories: department
            },

            yAxis: {
                categories: weekday,
                title: '요일'
            },

            colorAxis: {
                min: 0,
                minColor: '#FFFFFF',
                maxColor: Highcharts.getOptions().colors[0]
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
                    return '<b>' + this.series.xAxis.categories[this.point.x] + '</b> 수술 <br><b>' +
                        this.point.value + '</b>번<br><b>' + this.series.yAxis.categories[this.point.y] + '요일</b>';
                }
            },

            series: [{
                name: 'Sales per employee',
                borderWidth: 1,
                data: chartdata,
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
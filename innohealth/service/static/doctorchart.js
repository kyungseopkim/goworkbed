$(function(){$.getJSON("http://localhost:1323/operation/bydoctor?dbname=testdb",function(t){opsdata=[],doctors=[],$.each(t,function(t,e){doctors.push(e["의사"]),opsdata.push(e["수술"])}),Highcharts.chart("container",{chart:{type:"column"},title:{text:"Operation By Doctors"},xAxis:{categories:doctors,scrollbar:{enabled:!0},ticklength:0,title:{text:null},min:0,max:10},yAxis:{title:{text:"Operation"}},series:[{name:"수술",data:opsdata}],plotOptions:{bar:{dataLabels:{enabled:!0}}},legend:{enabled:!1},credits:{enabled:!1}}),Highcharts.chart("container2",{chart:{type:"column"},title:{text:"Operation By Doctors"},xAxis:{categories:doctors,ticklength:0,title:{text:null}},yAxis:{title:{text:"Operation"}},series:[{name:"수술",data:opsdata}],plotOptions:{bar:{dataLabels:{enabled:!0}}},legend:{enabled:!1},credits:{enabled:!1}})})})
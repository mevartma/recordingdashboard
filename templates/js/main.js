/**
 * Created by mevratm on 3/13/17.
 */

function GetDataByNumber() {
    var number = document.getElementById("number").value;
    var date1 = document.getElementById("date1").value;
    var date2 = document.getElementById("date2").value;
    var office = document.getElementById("office").value;
    console.log(number,date1,date2,office);
    var getData = {
        "command": "number",
        "number": number,
        "date1": date1,
        "date2": date2,
        "office": office
    };

    $.ajax({
        url: "/api/v1/recordings",
        type: "get",
        data: getData,
        success: function(data){
            $('#recordingdata').dataTable().fnDestroy();
            $("#recordingdata").dataTable({
                data: data,
                "columns": [
                    { "data": "calldate", "orderable" : true },
                    { "data": "cnam", "orderable" : true },
                    { "data": "src", "orderable" : true },
                    { "data": "dst", "orderable" : true },
                    { "data": "billsec" },
                    { "data": "disposition", "orderable" : true },
                    { "render": function(data,type,full,meta){
                        if(full.disposition == "NO ANSWER"){
                            return '<lable>NO File</lable>';
                        } else {
                            return '<audio controls preload="none"><source src="'+full.server_url+'" type="audio/mpeg"></audio>';
                        }
                    }},
                    { "data": "office", "orderable" : true },
                ]
            });
            /*$("#recordingdata thead th").each(function () {
                var title = $(this).text();
                $(this).html( '<input type="text" placeholder="Search '+title+'" />' );
            });

            var table = $("#recordingdata").dataTable();

            table.api().columns().every(function () {
                var that = this;
                $('input', this.footer()).on('keyup change', function () {
                    if (that.search() !== this.value ) {
                        that.search(this.value).draw();
                    }
                });
            });*/
        }
    });
}

/*function GetDataByNumber() {
 var number = document.getElementById("number").value;
 var formData = {
 "command": "number",
 "number": number.toString()
 };

 $("#recordingdata").dataTable({
 "ajax": {
 "url": "/api/v1/recordings",
 "data": formData,
 "dataSrc": ""
 },
 "columns": [
 { "data": "id", "orderable" : true },
 { "data": "calldate", "orderable" : true },
 { "data": "src", "orderable" : true },
 { "data": "dst", "orderable" : true },
 { "data": "duration" },
 { "data": "billsec" },
 { "data": "disposition", "orderable" : true },
 { "render": function(data,type,full,meta){
 if(full.disposition == "NO ANSWER"){
 return '<lable>NO File</lable>';
 } else {
 //return '<audio controls><source src="'+full.s_3_file_url+'" type="audio/wav" preload="none"></audio>';
 return '<a href="'+full.s_3_file_url+'">Download</a>';
 }
 }},
 { "data": "office", "orderable" : true },
 ]
 });

 $("#recordingdata thead th").each(function () {
 var title = $(this).text();
 $(this).html( '<input type="text" placeholder="Search '+title+'" />' );
 });

 var table = $("#recordingdata").dataTable();

 table.api().columns().every(function () {
 var that = this;
 $('input', this.footer()).on('keyup change', function () {
 if (that.search() !== this.value ) {
 that.search(this.value).draw();
 }
 });
 });
 }*/

/*$(document).ready(function (event) {
 var formData = {
 "command": "range",
 "from": "1",
 "to": "100"
 };

 $("#recordingdata").dataTable({
 "ajax": {
 "url": "/api/v1/recordings",
 "data": formData,
 "dataSrc": ""
 },
 "columns": [
 { "data": "id", "orderable" : true },
 { "data": "calldate", "orderable" : true },
 { "data": "src", "orderable" : true },
 { "data": "dst", "orderable" : true },
 { "data": "duration" },
 { "data": "billsec" },
 { "data": "disposition", "orderable" : true },
 { "render": function(data,type,full,meta){
 if(full.disposition == "NO ANSWER"){
 return '<lable>NO File</lable>';
 } else {
 //return '<audio controls><source src="'+full.s_3_file_url+'" type="audio/wav" preload="none"></audio>';
 return '<a href="'+full.s_3_file_url+'">Download</a>';
 }
 }},
 { "data": "office", "orderable" : true },
 ]
 });

 $("#recordingdata thead th").each(function () {
 var title = $(this).text();
 $(this).html( '<input type="text" placeholder="Search '+title+'" />' );
 });

 var table = $("#recordingdata").dataTable();

 table.api().columns().every(function () {
 var that = this;
 $('input', this.footer()).on('keyup change', function () {
 if (that.search() !== this.value ) {
 that.search(this.value).draw();
 }
 });
 });

 //event.preventDefault();
 });*/

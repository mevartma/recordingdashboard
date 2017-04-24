/**
 * Created by mevratm on 3/13/17.
 */

function GetDataByNumber() {
    var number = document.getElementById("number").value;
    console.log(number);
    var getData = {
        "command": "number",
        "number": number
    };

    $.ajax({
        url: "/api/v1/recordings",
        type: "get",
        data: getData,
        success: onSuccess(response),
        error: onError(xhr)
    });
}

function onSuccess(data) {
    console.log(data);
}

function onError(er) {
    console.log(er);
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
function sorting(json_object, key_to_sort_by)
{
    function sortByKey(a, b) {
        var x = a[key_to_sort_by];
        var y = b[key_to_sort_by];
        return ((x < y) ? -1 : ((x > y) ? 1 : 0));
    }

    json_object.sort(sortByKey);

    var size = values.length;
    var i=0;



}
$(document).ready(function()
{
    // init modal and selector
    $('.modal-trigger').leanModal();
    $('select').material_select();

    // get data table object
    table = $('#datatable').DataTable();
    document.getElementById("subHeader").innerHTML=0+" รายการ";



    // restore cookies ////////
    if (localStorage["fixCar"])
        RestoreCookies();

});

function checkExist(id)
{
    var size = values.length;
    var check =false;
    var i;
    var j;

    for(i = 0 ; i < size && check!=true ; i++)
    {
        if(values[i].id==id)
        {
            j=i;
            check = true;
        }
    }

    if(check==true)
    return j;
    else return-1;
}



$("#pick").click( function ()
{
    var temp;




    $('#datatable input[ name="check[]" ]:checked').each(function ()
    {
        row = $(this).closest("tr");
        item= ({

            id: $(row).find('input[name="check[]"]:checked').attr('id') ,
            name: table.row(row).data()[1] ,
            type: table.row(row).data()[2] ,
            amount: 1 ,
            unit: table.row(row).data()[3] ,
            price: table.row(row).data()[4]
        });


        id = item.id;
        index=checkExist(id);

        // รายการที่มีแล้ว ให้เพิ่มจำนวน
        if( index!=-1 )
        {

                // console.log("value size after add " + values.length);
                // console.log("index " + index);
                // console.log("values[index].amount " + values[index].amount);
                values[index].amount = parseInt(values[index].amount )+1;
                Materialize.toast('รายการ ' + values[index].name + ' มีจำนวน ' + values[index].amount + ' ชิ้น', 2000);


        }
        // รายการยังไม่มี ให้ใส่ในอาร์เรย์
        else
        {

            values.push(item);


            index = values.length-1;
            temp = values[index].unit.toString();
            values[index].unit = temp.substr( temp.search(" "),temp.length );

            // console.log("firstly add ");
            // console.log("values[index].amount " + values[index].amount);

            document.getElementById("badge").style.display = "inline";
            document.getElementById("badge").innerHTML=values.length;
            document.getElementById("subHeader").innerHTML=values.length+" รายการ";
            Materialize.toast('อะไหล่รถยนต์ ' + values[index].name + ' ถูกเพิ่มในรายการ', 2000);
            sorting(values, 'id');



        }




        $("#SparePart ").empty();
        localStorage.setItem("fixCar", JSON.stringify(values) );
        RestoreCookies();







    });



} );


$("#fix ").click( function ()
{

    $('#report').openModal
    (
        {

            complete: function()
            {

                $('#SparePart input[ name="amount[]" ]').each(function ()
                {
                    // alert($(this).closest("a").attr('value'));
                    console.log($(this).attr('id'));
                    console.log($(this).attr('value'));
                    checkExist( $(this).attr('id').toString().substring(3,$(this).attr('id').toString().length ) ) ;



                });
            }
        }
    );



} );




$("#done").click( function ()
{








} );



$("#SparePart").on( 'click', 'li', function ()
{




    // console.log( $(this).closest("li").attr('id') );
    console.log( "value size " +values.length  );
    // console.log( indexOfLI );




} );


// delete a row from collection
$("#SparePart").on( 'click', 'a', function ()
{
    indexOfLI = checkExist( $(this).closest("a").attr('value') ) ;

    if(indexOfLI!=-1)
    {
        values.splice(indexOfLI, 1);
        console.log("value size after delete " + values.length);
        $("#SparePart ").empty();
        localStorage.setItem("fixCar", JSON.stringify(values) );
        RestoreCookies();
    }


});




function RestoreCookies()
{

    values = JSON.parse(localStorage.getItem("fixCar"));
    var size = values.length;
    document.getElementById("badge").style.display = "inline";
    document.getElementById("badge").innerHTML = size;
    document.getElementById("subHeader").innerHTML = size + " รายการ";
    var i;
    var total=0;



    for(i = 0 ; i<size ; i++)
    {

        total = total+( parseFloat(values[i].price)*parseFloat( values[i].amount ) );

        $("#SparePart ").append('<li id="'+values[i].id+'" class="collection-item">'
            // +'<br>'
            +'<div class="row">'
            +'<div class="col s5">'
            +'<span>'+ values[i].name +'</span>'
            +'</div>'

            +'<div class="col s3">'
            +'<span>'+ parseFloat(values[i].price)+' บาท' +'</span>'
            +'</div>'

            // +'<div class="col s2">'
            // +'<span>'+'&emsp;'+ values[i].amount +'</span>'
            // +'<span>'+''+ values[i].unit +'</span>'
            // +'</div>'
            +'<div class="col s3">'
            +'<input min="0" max="10" style="width: 46px; height: 16px;" type="number" class="center checkbox-indigo filled-in" id="id+'+values[i].id+'" value="'+values[i].amount+'" name="amount[]"  /> <label for="id+'+values[i].id+'"></label>'
            +'<span>'+' / '+ values[i].unit +'</span>'
            +'</div>'

            +'<div class="col s1">'
            +'<a  class="right" type="button" value="'+values[i].id+'"><i class="red-text material-icons ">clear</i></a>'
            +'</div>'

            +'</div>'
            +'</li>');


    }


    $("#SparePart ").append('<li class="collection-item">'
        +'<br>'
        +'<div class="row">'
        +'<div class="col s5">'
        +'<span class="orange-text">'+ 'สรุปยอดการสั่งซ่อม'+'</span>'
        +'</div>'

        +'<div class="col s7">'
        +'<span class="orange-text">'+ parseFloat(total)+' บาท' +'</span>'
        +'</div>'
        +'</div>'

        +'</li>');



}








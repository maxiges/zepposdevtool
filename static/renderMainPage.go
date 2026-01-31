package static

var RenderPageTop = `
 <div class="container sticky-top bg-dark-subtle mt-0 w-100 rounded-5 rounded-top-0">
 <div class="row pt-2">
 <div class="col-md-auto align-self-center text-light">
 	<div class="bg-success rounded px-2">
 		<p class="text-center fw-bold"> App name <p>
		<p class="text-center"> {{.AppName}} <p>
		</div>
	</div>

	<div class="col-md-auto">
	<div class="bg-info rounded px-2 text-light">
		<p class="text-center"> Total memory </p>
		<p class="dropdown-divider"></p>
		<p class="text-center"> {{.Total}}  ->  {{.TotalKB}}KB</p>
		</div>
	</div>

	<div class="col-md-auto">
		<a href="javascript:window.location.href=window.location.href" class="mt-2 btn btn-primary btn-lg" tabindex="-1" role="button" >Reload ⟳</a>
	</div>

	<div class="col-md-auto"> 
		<div class="rounded px-2">
			<p>
				<label for="change_refresh">Auto Refresh</label>
			</p>
			 <select 
            name="change_refresh" 
            id="select-refresh-time" 
            onchange="autoRefresh(this);" 
            class="form-select form-select-sm w-auto" 
            style="min-width: 80px;"
        >
            <option value="0">Off</option>
            <option value="2000">2s</option>
            <option value="5000">5s</option>
            <option value="10000">10s</option>
            <option value="30000">30s</option>
            <option value="60000">60s</option>
        </select>
		</div>
	</div>
</div>
</div>



 <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz" crossorigin="anonymous"></script>
<script>

var timeout = null;
window.onload = function() {
    var refreshTime = sessionStorage.getItem("refreshTime");
    if (refreshTime != null && Number(refreshTime) > 0) {
	 document.getElementById("select-refresh-time").value = refreshTime;
	  clearTimeout(timeout);
        timeout = setTimeout(() => {
            window.location.reload();
        }, refreshTime);
    }
}

function autoRefresh(obj) {
    sessionStorage.setItem("refreshTime", obj.value);
    clearTimeout(timeout);
    if (obj.value == "0") {} else {
        timeout = setTimeout(() => {
            window.location.reload();
        }, obj.value);
    }
}
</script>
`

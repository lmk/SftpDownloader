/*
별도의 html 파일 없이 실행 파일에 포함하기 위한 코드
*/
package main

import "fmt"

const HTML_ROOT = `
<!DOCTYPE html>
<html lang="kr">
<head>
	<meta charset="utf-8">
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@4.0.0/dist/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
	<script src="https://code.jquery.com/jquery-3.2.1.slim.min.js" integrity="sha384-KJ3o2DKtIkvYIK3UENzmM7KCkRr/rE9/Qpg6aAZGJwFDMVNA/GpGFF93hXpG5KkN" crossorigin="anonymous"></script>
	<script src="https://cdn.jsdelivr.net/npm/popper.js@1.12.9/dist/umd/popper.min.js" integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous"></script>
	<script src="https://cdn.jsdelivr.net/npm/bootstrap@4.0.0/dist/js/bootstrap.min.js" integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>  
	<title>Sftp Downloader</title>
</head>
<body class="bg-light">
	<div class="container">
		<div class="py-5 text-center">
			<h2>Sftp Downloader</h2>
			<p class="lead">경로를 유지하고 파일 전체를 다운로드 받습니다</p>
		</div>
		<div class="col-md-12">
			<form action="http://localhost:%s/download" method="post">
				<h4 class="mb-3">SFTP Server</h4>
				<div class="row mb-1">
					<div class="input-group col-md-6 mb-3">
						<label class="input-group-text" for="sftp-addr">IP Address</label>
						<input type="text" class="form-control" name="sftp-addr" placeholder="" value="%s" required>
					</div>
				</div>
				<div class="row mb-3">
					<div class="input-group col-md-3 mb-3">
						<label class="input-group-text" for="sftp-id">ID</label>
						<input type="text" class="form-control" name="sftp-id" placeholder="" value="%s" required>
					</div>
					<div class="input-group col-md-3 mb-3">
						<label class="input-group-text" for="sftp-pwd">Password</label>
						<input type="text" class="form-control" name="sftp-pwd" placeholder="" value="%s" required>
					</div>
				</div>
				<h4 class="mb-3">Local</h4>
				<div class="row mb-3">
					<div class="input-group col-md-12 mb-3">
						<label class="input-group-text" for="local-dir">Local Directory</label>
						<input type="text" class="form-control" name="local-dir" placeholder="" value="%s" required>
					</div>
				</div>
				<h4 class="mb-3">File List</h4>
				<div class="row mb-3">
					<div class="col-md-12">
						<textarea class="col-md-12" rows="10" name="file-list" placeholder="File List...">%s</textarea>
					</div>
				</div>
				<hr class="mb-3">
				<button class="btn btn-primary btn-lg btn-block" type="submit">Download</button>
			</form>
		</div>
	</div>
</body>
</html>
`

const HTML_DOWNLOAD = `
<!DOCTYPE html>
<html lang="kr">
<head>
    <meta charset="utf-8">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@4.0.0/dist/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
    <script src="https://code.jquery.com/jquery-3.2.1.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/popper.js@1.12.9/dist/umd/popper.min.js" integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous"></script>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@4.0.0/dist/js/bootstrap.min.js" integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>  
    <title>Sftp Downloader</title>
</head>
<body class="bg-light">
    <div class="container">
        <div class="py-5 text-center">
            <h2>Sftp Downloader</h2>
            <p class="lead">경로를 유지하고 파일 전체를 다운로드 받습니다</p>
        </div>
        <div class="col-md-12">
            <table class="table">
                <thead>
                    <tr>
                    <th scope="col">#</th>
                    <th scope="col">Stat</th>
                    <th scope="col">Path</th>
                    <th scope="col">Date</th>
                    <th scope="col">Size</th>
                    </tr>
                </thead>
                <tbody>
                    %s
                </tbody>
                </table>
        </div>
    </div>
	<style>
		td.path {
			position: relative;
		}
		td.path div.progress{
			position: absolute;
			top: 0;
			left: 0;
			width: 0%%;
			height: 100%%;
			background-color: lightblue;
			z-index: -1;
		}
	</style>
	<script>
		function humanFileSize(bytes, si=false, dp=1) {
			const thresh = si ? 1000 : 1024;
		
			if (Math.abs(bytes) < thresh) {
			return bytes + ' B';
			}
		
			const units = si 
			? ['kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'] 
			: ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
			let u = -1;
			const r = 10**dp;
		
			do {
			bytes /= thresh;
			++u;
			} while (Math.round(Math.abs(bytes) * r) / r >= thresh && u < units.length - 1);
		
		
			return bytes.toFixed(dp) + ' ' + units[u];
		}


		var x = setInterval(function() {
			$.ajax({
				type: 'post',
				url: 'http://localhost:%s/downloading',
				success : function(result) {

					result.files.forEach((file, index)=>{

						clocks = ["🕛", "🕐","🕑","🕒","🕓","🕔","🕕","🕖","🕗","🕘", "🕙", "🕚" ]

						objDanger = $('[id="'+file.path+'"]').children('td[name="path"]').children(".text-danger")
			
						if (typeof file.notify == "undefined" || file.notify=="") {
							ch = ""
							idx = file.localSize / file.remoteSize * (clocks.length-1)
							per = Math.round(file.localSize / file.remoteSize * 100)
			
							if (file.remoteSize == 0) {
								ch = '❌'
							} else if (isNaN(idx)) {
								ch = clocks[0]
							} else if (file.localSize == file.remoteSize) {
								ch = '✔'
							} else {
								ch = clocks[Math.round(idx)]
							}
			
							$('[id="'+file.path+'"]').children('td[name="stat"]').text(ch)
							$('[id="'+file.path+'"]').children('td[name="date"]').text(file.date)
							$('[id="'+file.path+'"]').children('td[name="size"]').text(humanFileSize(file.localSize))
			
							if (file.localSize == file.remoteSize) {
								$('[id="'+file.path+'"]').children('td[name="path"]').children(".progress").css("width", "0%%")
							} else {
								$('[id="'+file.path+'"]').children('td[name="path"]').children(".progress").css("width", Math.round(file.localSize / file.remoteSize * 100) + "%%")
							}

							$(objDanger).remove();
			
						} else {
							if (objDanger.length == 0) {
								html = '<div class="text-danger">'+file.notify+'</div>'
								$('[id="'+file.path+'"]').children('td[name="path"]').append(html)
							} else {
								$(objDanger).children("div").text(file.notify)
							}
						}
					})

					if (result.stat == "READY") {
						clearInterval(x)
						return
					} 

					if (result.stat == "DONE") {
						clearInterval(x)
						setTimeout(() => alert('Download complete!\nPlease, Close this window.'), 1000);
						return
					} 

				},
				error : function(xhr, status, message) {
					clearInterval(x)
					console.log("error : "+message);
					window.open('','_self').close(); 
				}
			})
		}, 100);
	</script>
</body>
</html>
`

const HTML_DOWNLOAD_ROW = `
<tr id="%s">
	<th scope="row">%d</th>
	<td name="stat">%s</td>
	<td name="path" class="path"><div class="progress"></div><div>%s</div></td>
	<td name="date">%s</td>
	<td name="size">0B</td>
</tr>
`

func HtmlRoot() string {
	return fmt.Sprintf(HTML_ROOT,
		AppPort,
		downInfo.Ip,
		downInfo.Id,
		downInfo.Password,
		downInfo.LocalDir,
		downInfo.RemoteFilesString())
}

func HtmlDownload() string {

	html := ""
	for i, file := range downInfo.Files {
		if file.Remote.IsExist {
			html += fmt.Sprintf(HTML_DOWNLOAD_ROW, file.Remote.Path, i+1, "✔", file.Remote.Path, file.Remote.Date)
		} else {
			html += fmt.Sprintf(HTML_DOWNLOAD_ROW, file.Remote.Path, i+1, "❌", file.Remote.Path, file.Remote.Date)
		}
	}

	return fmt.Sprintf(HTML_DOWNLOAD, html, AppPort)
}

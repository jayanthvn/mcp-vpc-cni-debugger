#!/usr/bin/env python3
"""
ALB MCP Server - Proper MCP Protocol Implementation
Provides detailed Kubernetes pod network information via MCP protocol
"""

import json
import sys
import requests
import logging
from typing import Dict, Any, List

# Configure logging to stderr so it doesn't interfere with MCP communication
logging.basicConfig(level=logging.INFO, stream=sys.stderr, 
                   format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

ALB_BASE_URL = "add ALB url"

class ALBMCPServer:
    def __init__(self):
        self.tools = [
            {
                "name": "get_pod_network_info",
                "description": "Get detailed network information for a Kubernetes pod including ENI details, security groups, and routing rules",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "namespace": {
                            "type": "string",
                            "description": "Kubernetes namespace where the pod is located"
                        },
                        "pod_name": {
                            "type": "string", 
                            "description": "Name of the pod to get network information for"
                        }
                    },
                    "required": ["namespace", "pod_name"]
                }
            }
        ]
    
    def get_pod_network_info(self, namespace: str, pod_name: str) -> Dict[str, Any]:
        """Get network information for a specific pod from ALB endpoint"""
        url = f"{ALB_BASE_URL}/mcp/network/pod/{namespace}/{pod_name}"
        try:
            logger.info(f"Querying ALB endpoint: {url}")
            response = requests.get(url, timeout=30)
            response.raise_for_status()
            data = response.json()
            logger.info(f"Successfully retrieved network info for pod {namespace}/{pod_name}")
            return data
        except requests.exceptions.RequestException as e:
            error_msg = f"Failed to query ALB MCP server: {str(e)}"
            logger.error(error_msg)
            return {"error": error_msg}
    
    def handle_initialize(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Handle MCP initialize request"""
        return {
            "protocolVersion": "2024-11-05",
            "capabilities": {
                "tools": {}
            },
            "serverInfo": {
                "name": "alb-mcp-server",
                "version": "1.0.0"
            }
        }
    
    def handle_list_tools(self) -> Dict[str, Any]:
        """Handle MCP tools/list request"""
        return {"tools": self.tools}
    
    def handle_call_tool(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Handle MCP tools/call request"""
        tool_name = params.get("name")
        arguments = params.get("arguments", {})
        
        if tool_name == "get_pod_network_info":
            namespace = arguments.get("namespace")
            pod_name = arguments.get("pod_name")
            
            if not namespace or not pod_name:
                return {
                    "content": [
                        {
                            "type": "text",
                            "text": "Error: Both 'namespace' and 'pod_name' parameters are required"
                        }
                    ],
                    "isError": True
                }
            
            result = self.get_pod_network_info(namespace, pod_name)
            
            if "error" in result:
                return {
                    "content": [
                        {
                            "type": "text", 
                            "text": f"Error: {result['error']}"
                        }
                    ],
                    "isError": True
                }
            else:
                # Format the network information nicely
                formatted_result = self.format_network_info(result)
                return {
                    "content": [
                        {
                            "type": "text",
                            "text": formatted_result
                        }
                    ]
                }
        else:
            return {
                "content": [
                    {
                        "type": "text",
                        "text": f"Error: Unknown tool '{tool_name}'"
                    }
                ],
                "isError": True
            }
    
    def format_network_info(self, data: Dict[str, Any]) -> str:
        """Format network information for readable output"""
        result = []
        result.append(f"Pod Network Information for {data.get('namespace', 'N/A')}/{data.get('podName', 'N/A')}")
        result.append("=" * 60)
        result.append(f"Pod IP: {data.get('podIP', 'N/A')}")
        
        if 'eni' in data:
            eni = data['eni']
            result.append("\nENI Details:")
            result.append(f"  ENI ID: {eni.get('eniId', 'N/A')}")
            result.append(f"  Device: {eni.get('device', 'N/A')}")
            result.append(f"  MAC Address: {eni.get('mac', 'N/A')}")
            result.append(f"  Subnet: {eni.get('subnet', 'N/A')}")
            result.append(f"  VPC: {eni.get('vpc', 'N/A')}")
            
            if 'sgIds' in eni:
                result.append(f"  Security Groups: {', '.join(eni['sgIds'])}")
        
        if 'routeRules' in data:
            result.append(f"\nRouting Rules ({len(data['routeRules'])} rules):")
            for i, rule in enumerate(data['routeRules'][:5], 1):  # Show first 5 rules
                result.append(f"  {i}. {rule.get('rule', 'N/A')}")
            
            if len(data['routeRules']) > 5:
                result.append(f"  ... and {len(data['routeRules']) - 5} more rules")
        
        return "\n".join(result)
    
    def process_request(self, request: Dict[str, Any]) -> Dict[str, Any]:
        """Process incoming MCP request"""
        method = request.get("method")
        params = request.get("params", {})
        
        logger.info(f"Processing MCP request: {method}")
        
        if method == "initialize":
            return self.handle_initialize(params)
        elif method == "tools/list":
            return self.handle_list_tools()
        elif method == "tools/call":
            return self.handle_call_tool(params)
        else:
            return {
                "error": {
                    "code": -32601,
                    "message": f"Method not found: {method}"
                }
            }
    
    def run(self):
        """Main server loop - reads JSON-RPC requests from stdin and writes responses to stdout"""
        logger.info("ALB MCP Server starting...")
        
        try:
            for line in sys.stdin:
                line = line.strip()
                if not line:
                    continue
                
                try:
                    request = json.loads(line)
                    logger.info(f"Received request: {request.get('method', 'unknown')}")
                    
                    response = self.process_request(request)
                    
                    # Add request ID to response if present
                    if "id" in request:
                        response["id"] = request["id"]
                    
                    # Write response to stdout
                    print(json.dumps(response), flush=True)
                    logger.info(f"Sent response for: {request.get('method', 'unknown')}")
                    
                except json.JSONDecodeError as e:
                    logger.error(f"Invalid JSON received: {e}")
                    error_response = {
                        "error": {
                            "code": -32700,
                            "message": "Parse error"
                        }
                    }
                    print(json.dumps(error_response), flush=True)
                except Exception as e:
                    logger.error(f"Error processing request: {e}")
                    error_response = {
                        "error": {
                            "code": -32603,
                            "message": f"Internal error: {str(e)}"
                        }
                    }
                    if "id" in request:
                        error_response["id"] = request["id"]
                    print(json.dumps(error_response), flush=True)
                    
        except KeyboardInterrupt:
            logger.info("ALB MCP Server shutting down...")
        except Exception as e:
            logger.error(f"Fatal error: {e}")
            sys.exit(1)

def main():
    """Entry point for the MCP server"""
    server = ALBMCPServer()
    server.run()

if __name__ == "__main__":
    main()

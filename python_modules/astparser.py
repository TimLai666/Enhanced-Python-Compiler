# astparser.py
import ast
import json

def parse_code(source_code: str) -> str:
    parsed_ast = ast.parse(source_code)
    
    def ast_to_json(node):
        if isinstance(node, ast.AST):
            fields = {field: ast_to_json(getattr(node, field)) for field in node._fields}
            return {'_type': node.__class__.__name__, **fields}
        elif isinstance(node, list):
            return [ast_to_json(item) for item in node]
        else:
            return node

    return json.dumps(ast_to_json(parsed_ast))
